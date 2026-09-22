provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_cloudcontrol_guardrail" "guardrail" {
  guardrail_id    = var.guardrail_id
  landing_zone_id = var.landing_zone_id != "" ? var.landing_zone_id : null
}