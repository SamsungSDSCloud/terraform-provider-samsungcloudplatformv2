provider "samsungcloudplatformv2" {
}

resource "samsungcloudplatformv2_cloudcontrol_guardrail_binding" "guardrail_binding" {
  guardrail_ids   = var.guardrail_ids
  unit_ids        = var.unit_ids
  landing_zone_id = var.landing_zone_id != "" ? var.landing_zone_id : null
}