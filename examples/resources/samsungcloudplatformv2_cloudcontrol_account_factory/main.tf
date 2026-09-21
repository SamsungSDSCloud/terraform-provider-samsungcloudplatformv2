provider "samsungcloudplatformv2" {
}

resource "samsungcloudplatformv2_cloudcontrol_account_factory" "account_factory" {
  landing_zone_id   = var.landing_zone_id
  name             = var.name
  parent_unit_id   = var.parent_unit_id
  email            = var.email
  sso_user_email   = var.sso_user_email
  sso_user_name    = var.sso_user_name
  sso_user_real_name = var.sso_user_real_name
}