package basicldapauth

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/require"
	"gotest.tools/v3/assert"
	"io"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
)

type VaultLdapConfig struct {
	RequestID     string `json:"request_id"`
	LeaseID       string `json:"lease_id"`
	Renewable     bool   `json:"renewable"`
	LeaseDuration int    `json:"lease_duration"`
	Data          struct {
		AnonymousGroupSearch     bool          `json:"anonymous_group_search"`
		Binddn                   string        `json:"binddn"`
		CaseSensitiveNames       bool          `json:"case_sensitive_names"`
		Certificate              string        `json:"certificate"`
		DenyNullBind             bool          `json:"deny_null_bind"`
		Discoverdn               bool          `json:"discoverdn"`
		Groupattr                string        `json:"groupattr"`
		Groupdn                  string        `json:"groupdn"`
		Groupfilter              string        `json:"groupfilter"`
		InsecureTlS              bool          `json:"insecure_tls"`
		RequestTimeout           int           `json:"request_timeout"`
		Starttls                 bool          `json:"starttls"`
		TLSMaxVersion            string        `json:"tls_max_version"`
		TLSMinVersion            string        `json:"tls_min_version"`
		TokenBoundCidrs          []interface{} `json:"token_bound_cidrs"`
		TokenExplicitMaxTTL      int           `json:"token_explicit_max_ttl"`
		TokenMaxTTL              int           `json:"token_max_ttl"`
		TokenNoDefaultPolicy     bool          `json:"token_no_default_policy"`
		TokenNumUses             int           `json:"token_num_uses"`
		TokenPeriod              int           `json:"token_period"`
		TokenPolicies            []interface{} `json:"token_policies"`
		TokenTTL                 int           `json:"token_ttl"`
		TokenType                string        `json:"token_type"`
		Upndomain                string        `json:"upndomain"`
		URL                      string        `json:"url"`
		UsePre111GroupCnBehavior bool          `json:"use_pre111_group_cn_behavior"`
		UseTokenGroups           bool          `json:"use_token_groups"`
		Userattr                 string        `json:"userattr"`
		Userdn                   string        `json:"userdn"`
		Userfilter               string        `json:"userfilter"`
		UsernameAsAlias          bool          `json:"username_as_alias"`
	} `json:"data"`
	WrapInfo interface{} `json:"wrap_info"`
	Warnings interface{} `json:"warnings"`
	Auth     interface{} `json:"auth"`
}

type VaultLdapGroup struct {
	RequestID     string `json:"request_id"`
	LeaseID       string `json:"lease_id"`
	Renewable     bool   `json:"renewable"`
	LeaseDuration int    `json:"lease_duration"`
	Data          struct {
		Policies []string `json:"policies"`
	} `json:"data"`
	WrapInfo interface{} `json:"wrap_info"`
	Warnings interface{} `json:"warnings"`
	Auth     interface{} `json:"auth"`
}

func TestBasicLdapAuthEndpoint(t *testing.T) {
	t.Parallel()

	fmt.Println("Executing")

	vaultToken := os.Getenv("VAULT_TOKEN")
	if vaultToken == "" {
		fmt.Println("Did not find VAULT_TOKEN environment variable, reverting to user file token.")
		currentUser, err := user.Current()
		if err != nil {
			fmt.Println(err.Error())
			t.Error(err.Error())
		}
		fullpath, err := filepath.Abs(fmt.Sprintf("/Users/%s/.vault-token", currentUser.Username))
		if err != nil {
			fmt.Println(err.Error())
			t.Fatalf("Unable to resolve absolute path of file")
		}
		fileBytes, err := os.ReadFile(fullpath)
		if err != nil {
			fmt.Println(err.Error())
			t.Error(err.Error())
		} else {
			fmt.Println(fmt.Sprintf("Found User: %s Vault Token", currentUser.Username))
			vaultToken = string(fileBytes)
		}
	}
	vaultAddr := os.Getenv("VAULT_ADDR")
	fmt.Println(fmt.Sprintf("Using Vault Address: %s", vaultAddr))
	moduleVars := map[string]interface{}{
		"vault_address": vaultAddr,
	}
	fmt.Println("Module Vars: ", moduleVars)

	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../../examples/basic_ldap_auth",
		Vars:         moduleVars,
	})

	// Cleanup resources after end of test with "terraform destroy"
	defer terraform.Destroy(t, terraformOptions)

	// Run "terraform init" and "terraform apply". Fail the test if there are any errors.
	terraform.InitAndApply(t, terraformOptions)

	fmt.Println("Testing API")

	// Inspect Vault LDAP Auth Config
	url := fmt.Sprintf("%s/v1/auth/__test/ldap/config", vaultAddr)
	body := GetAPI(url, vaultToken)
	//fmt.Println(string(body))

	var config VaultLdapConfig
	err := json.Unmarshal(body, &config)
	if err != nil {
		fmt.Println(err.Error())
	}

	assert.Equal(t, config.Data.Binddn, "uid=ThisIsNotTheBindUidYouAreLookingFor,ou=Users,o=ThisIsNotTheAdOrgIdYouAreLookingFor,dc=jumpcloud,dc=com")
	require.Contains(t, config.Data.Certificate, "-----BEGIN CERTIFICATE-----")
	assert.Equal(t, config.Data.Groupdn, "ou=Users,o=ThisIsNotTheAdOrgIdYouAreLookingFor,dc=jumpcloud,dc=com")
	assert.Equal(t, config.Data.Groupfilter, "(member={{.UserDN}})")
	assert.Equal(t, config.Data.Userdn, "ou=Users,o=ThisIsNotTheAdOrgIdYouAreLookingFor,dc=jumpcloud,dc=com")
	assert.Equal(t, config.Data.Userfilter, "({{.UserAttr}}={{.Username}})")
	assert.Equal(t, config.Data.URL, "ldaps://ldap.jumpcloud.com")

	// Inspect Vault LDAP Engineering Group
	url = fmt.Sprintf("%s/v1/auth/__test/ldap/groups/__test_engineering", vaultAddr)
	body = GetAPI(url, vaultToken)
	//fmt.Println(string(body))
	var group1 VaultLdapGroup
	err = json.Unmarshal(body, &group1)
	if err != nil {
		fmt.Println(err.Error())
	}
	assert.DeepEqual(t, group1.Data.Policies, []string{"__test_vault-identity-mgmt", "__test_vault-ops"})

	// Inspect Vault LDAP QA Group
	url = fmt.Sprintf("%s/v1/auth/__test/ldap/groups/__test_qa", vaultAddr)
	body = GetAPI(url, vaultToken)
	//fmt.Println(string(body))
	var group2 VaultLdapGroup
	err = json.Unmarshal(body, &group2)
	if err != nil {
		fmt.Println(err.Error())
	}
	assert.DeepEqual(t, group2.Data.Policies, []string{"__test_vault-qa"})
}

// Helper method to get Vault API for validation
func GetAPI(url string, vaultToken string) []byte {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println(err.Error())
	}
	req.Header.Add("X-Vault-Token", vaultToken)

	//res, err := http.DefaultClient.Do(req)
	// Skip Verify on API calls in tests
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	res, err := client.Do(req)
	if err != nil {
		fmt.Print(err.Error())
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("Failed to close API request handler!")
		}
	}(res.Body)
	body, readErr := io.ReadAll(res.Body)
	if readErr != nil {
		fmt.Print(err.Error())
	}
	//fmt.Println(string(body))

	return body
}
