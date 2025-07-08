provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_cloudmonitoring_event_notification_states" "eventNotificationStates" {
  x_resource_type = var.XResourceType
  event_id = var.EventId

}
