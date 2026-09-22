provider "samsungcloudplatformv2" {
}

resource "samsungcloudplatformv2_cloudcontrol_baseline_assignment" "baseline_assignment" {
  assignment_id      = var.assignment_id
  landing_zone_id    = var.landing_zone_id
  resource_type      = "OU"
  agree_yn           = var.agree_yn
  reregister_trigger = var.reregister_trigger
}
