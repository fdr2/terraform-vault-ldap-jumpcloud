locals {
  test_prefix = "__test"
}

module "vault-ldap-jumpcloud" {
  source                 = "../../"
  path                   = "${local.test_prefix}/ldap"
  jumpcloud_ad_bindpass  = var.vault_jumpcloud_ad_bindpass
  jumpcloud_ad_binduid   = var.vault_jumpcloud_ad_binduid
  jumpcloud_ad_orgid     = var.vault_jumpcloud_ad_orgid
  jumpcloud_group_policy = var.jumpcloud_group_policy
}

resource "vault_policy" "vault-ops" {
  name   = "${local.test_prefix}/vault-ops"
  policy = file("${path.root}/policy/vault-ops-policy.hcl")
}

resource "vault_policy" "vault-id-mgmt" {
  name   = "${local.test_prefix}/vault-identity-mgmt"
  policy = file("${path.root}/policy/vault-identity-mgmt-policy.hcl")
}

resource "vault_policy" "vault-qa" {
  name   = "${local.test_prefix}/vault-qa"
  policy = file("${path.root}/policy/vault-qa-policy.hcl")
}