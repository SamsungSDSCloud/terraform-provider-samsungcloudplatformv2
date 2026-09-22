provider "samsungcloudplatformv2" {
}

resource "samsungcloudplatformv2_cloudcontrol_baseline_assignment" "baseline_assignment" {
  assignment_id      = var.assignment_id
  landing_zone_id    = var.landing_zone_id
  resource_type      = "ACCOUNT"
  parent_unit_id     = var.parent_unit_id
  sso_user_name      = var.sso_user_name
  sso_user_real_name = var.sso_user_real_name
}