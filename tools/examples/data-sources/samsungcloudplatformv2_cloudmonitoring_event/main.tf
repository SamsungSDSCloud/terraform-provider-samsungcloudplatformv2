provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_cloudmonitoring_event" "eventDetail" {
  event_id = var.EventId
  x_resource_type = var.XResourceType

}
