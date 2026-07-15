provider "samsungcloudplatformv2" {
}

resource "samsungcloudplatformv2_organization_account" "account" {
  login_id        = var.login_id
  name            = var.name
  organization_id = var.organization_id
  role_name       = var.role_name
  lazy_policy     = var.lazy_policy
  parent_unit_id  = var.parent_unit_id
}