provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_cloudmonitoring_metrics" "metrics" {
  product_type_code = var.ProductTypeCode
#   object_type = var.ObjectType
}
