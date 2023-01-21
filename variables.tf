variable "path" {
  default = "ldap"
}

variable "jumpcloud_ad_binduid" {
  type = string
}

variable "jumpcloud_ad_bindpass" {
  type = string
}

variable "jumpcloud_ad_orgid" {
  type = string
}

variable "jumpcloud_ad_cert_filename" {
  type    = string
  default = null
}

variable "jumpcloud_group_policy" {
  type = map(object({
    policies : list(string)
    identities : map(object({
      disabled : bool
    }))
  }))
  #default = {
  #  Engineering : {
  #    policies : ["ops-vault-policy"]
  #    identities : {
  #      ops : {
  #        enabled : false
  #      }
  #    }
  #  }
  #}
}
