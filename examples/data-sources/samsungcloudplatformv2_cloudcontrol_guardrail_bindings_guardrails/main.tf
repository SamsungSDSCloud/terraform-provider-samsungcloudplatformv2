provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_cloudcontrol_guardrail_bindings_guardrails" "guardrails" {
  target_id       = var.target_id
  landing_zone_id = var.landing_zone_id != "" ? var.landing_zone_id : null
  name            = var.name != "" ? var.name : null
  size            = var.size > 0 ? var.size : null
  page            = var.page > 0 ? var.page : null
  sort            = var.sort != "" ? var.sort : null
}