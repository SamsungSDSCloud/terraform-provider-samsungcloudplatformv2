provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_cloudcontrol_guardrails" "guardrails" {
  landing_zone_id = var.landing_zone_id != "" ? var.landing_zone_id : null
  name            = var.name != "" ? var.name : null
  exclude_unit_id = var.exclude_unit_id != "" ? var.exclude_unit_id : null
  guidance        = var.guidance != "" ? var.guidance : null
  service_name    = var.service_name != "" ? var.service_name : null
  status          = var.status != "" ? var.status : null
  size            = var.size > 0 ? var.size : null
  page            = var.page > 0 ? var.page : null
  sort            = var.sort != "" ? var.sort : null
}