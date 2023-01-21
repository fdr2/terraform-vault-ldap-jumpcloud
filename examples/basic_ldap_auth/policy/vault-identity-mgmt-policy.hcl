# To create an entity and entity alias. Enable and configure Vault as an OIDC provider
path "identity/*" {
  capabilities = [ "create", "read", "update", "delete", "list" ]
}

# To enable userpass auth method
path "sys/auth/*" {
  capabilities = [ "create", "read", "update", "delete" ]
}

# To create a new user, "end-user" for userpass
path "auth/*" {
  capabilities = [ "create", "read", "update", "delete", "list" ]
}
