# Terraform Vault LDAP Jumpcloud

Configure Vault's LDAP authentication backend with Jumpcloud.

Issues a Vault token with the attached Nomad policy.

> This module aligns with [LDAP Auth Method](https://developer.hashicorp.com/vault/docs/auth/ldap)

## Usage
Add a module and assign a nomad policy for the tokens that will be issued.

```terraform
variable "jumpcloud_group_policy" {
  default = {
    Engineering : {
      policies : ["vault-ops-policy", "vault-identity-mgmt-policy"]
      identities : {
        ops : { disabled : false }
      }
    }
    QA : {
      policies : ["vault-ops-policy", ]
      identities : {
        qa : { disabled : true }
      }
    }
  }
}

module "vault-ldap-jumpcloud" {
  source = "./modules/terraform-vault-ldap-jumpcloud"
  jumpcloud_ad_bindpass  = var.vault_jumpcloud_ad_bindpass
  jumpcloud_ad_binduid   = var.vault_jumpcloud_ad_bindpass
  jumpcloud_ad_orgid     = var.vault_jumpcloud_ad_orgid
  jumpcloud_group_policy = var.jumpcloud_group_policy
}
```

## TODO

* write local vault dev example
* write tests
