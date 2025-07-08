provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_cloudmonitoring_event_policy" "event_policy" {
  x_resource_type = var.XResourceType
  event_policy_id = var.EventPolicyId
}
