provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_cloudmonitoring_event_policy_notifications" "eventPolicyNotifications" {
  event_policy_id = var.EventPolicyId
  x_resource_type = var.XResourceType

}
