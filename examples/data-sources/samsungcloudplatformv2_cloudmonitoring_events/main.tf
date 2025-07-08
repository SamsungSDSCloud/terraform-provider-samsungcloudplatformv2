provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_cloudmonitoring_events" "events" {
  event_state = var.EventState
  query_start_dt = var.QueryStartDt
  query_end_dt = var.QueryEndDt
  x_resource_type = var.XResourceType
  product_resource_id = var.ProductResourceId
}
