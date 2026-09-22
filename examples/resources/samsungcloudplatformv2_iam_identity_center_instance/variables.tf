variable "name" {
  type    = string
  default = "my-identity-center-instance2"
}

variable "description" {
  type    = string
  default = "My Identity Center Instance Description2"
}

variable "identity_store_type" {
  type    = string
  default = "IDENTITY_CENTER_DIRECTORY"
}

variable "identity_store_config" {
  type = object({
    bind_credential         = string
    bind_dn                 = string
    connection_url          = string
    rdn_ldap_attribute      = string
    user_object_classes     = string
    username_ldap_attribute = string
    users_dn                = string
  })
  default = null
}

variable "self_managed_password" {
  type    = bool
  default = false
}

variable "region" {
  type    = string
  default = "kr-west1"
}



