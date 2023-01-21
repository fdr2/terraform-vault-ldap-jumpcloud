variable "vault_address" {
  type    = string
  default = "http://127.0.0.1:8200"
}

variable "vault_jumpcloud_ad_bindpass" {
  default = "ThisIsNotTheBindPassYouAreLookingFor"
}

variable "vault_jumpcloud_ad_binduid" {
  default = "ThisIsNotTheBindUidYouAreLookingFor"
}

variable "vault_jumpcloud_ad_orgid" {
  default = "ThisIsNotTheAdOrgIdYouAreLookingFor"
}

variable "jumpcloud_group_policy" {
  default = {
    __test_Engineering : {
      policies : ["__test_vault-ops", "__test_vault-identity-mgmt"]
      identities : {
        __test_ops : { disabled : false }
      }
    }
    __test_QA : {
      policies : ["__test_vault-qa"]
      identities : {
        __test_qa : { disabled : true }
      }
    }
  }
}
