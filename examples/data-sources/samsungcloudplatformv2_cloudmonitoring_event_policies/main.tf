provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_cloudmonitoring_event_policies" "event_policies" {
  x_resource_type = var.XResourceType
  product_resource_id = var.ProductResourceId
#  metric_key = "c94607c8-6a15-4e13-84f3-99cb77d478e3"
  size = var.Size
  page = var.Page
}
