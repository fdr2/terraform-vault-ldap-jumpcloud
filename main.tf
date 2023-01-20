## Jumpcloud LDAP
locals {
  identity_names = flatten([for k, v in var.jumpcloud_group_policy : [for i, k in v.identities : i]])
  identity_map = merge([for i, g in var.jumpcloud_group_policy : {
    for k, v in g.identities : k => merge(v, { group : i })
  }]...)
}

resource "vault_ldap_auth_backend" "this" {
  path           = var.path
  url            = "ldaps://ldap.jumpcloud.com"
  certificate    = file(var.jumpcloud_ad_cert_filename)
  binddn         = "uid=${var.jumpcloud_ad_binduid},ou=Users,o=${var.jumpcloud_ad_orgid},dc=jumpcloud,dc=com"
  bindpass       = var.jumpcloud_ad_bindpass
  discoverdn     = false
  userdn         = "ou=Users,o=${var.jumpcloud_ad_orgid},dc=jumpcloud,dc=com"
  userattr       = "uid"
  groupdn        = "ou=Users,o=${var.jumpcloud_ad_orgid},dc=jumpcloud,dc=com"
  groupfilter    = "(member={{.UserDN}})"
  groupattr      = "memberOf"
  insecure_tls   = false
  starttls       = false
  deny_null_bind = true
}

resource "vault_ldap_auth_backend_group" "this" {
  for_each  = var.jumpcloud_group_policy
  groupname = each.key
  policies  = each.value["policies"]
  backend   = vault_ldap_auth_backend.this.path
}

resource "vault_identity_entity" "this" {
  for_each = local.identity_map
  name     = lower(each.key)
  disabled = try(each.value["disabled"], false)
}

resource "vault_identity_group" "this" {
  for_each          = var.jumpcloud_group_policy
  name              = lower(each.key)
  type              = "external"
  member_entity_ids = [for k, v in each.value["identities"] : k]
}

resource "vault_identity_entity_alias" "this" {
  for_each       = local.identity_map
  canonical_id   = vault_identity_entity.this[each.key].id
  mount_accessor = vault_ldap_auth_backend.this.accessor
  name           = lower(each.key)
}

resource "vault_identity_group_alias" "this" {
  for_each       = var.jumpcloud_group_policy
  canonical_id   = vault_identity_group.this[each.key].id
  mount_accessor = vault_ldap_auth_backend.this.accessor
  name           = lower(each.key)
}

resource "vault_identity_oidc_assignment" "this" {
  for_each = var.jumpcloud_group_policy
  name     = lower(each.key)
  entity_ids = concat(
    #[for k, v in each.value["identities"] : vault_identity_entity.this[k].id],
    [for k, v in each.value["identities"] : vault_identity_entity_alias.this[k].id],
  )
  group_ids = [
    #vault_identity_group.this[each.key].id,
    vault_identity_group_alias.this[each.key].id,
  ]
}
