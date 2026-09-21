provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_cloudcontrol_baseline_assignments" "baseline_assignments" {
  landing_zone_id = var.landing_zone_id != "" ? var.landing_zone_id : null
  resource_type   = var.resource_type != "" ? var.resource_type : null
  assignment_id   = var.assignment_id != "" ? var.assignment_id : null
  status          = var.status != "" ? var.status : null
}