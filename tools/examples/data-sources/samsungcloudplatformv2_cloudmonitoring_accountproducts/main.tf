provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_cloudmonitoring_accountproducts" "accountproducts" {
  x_resource_type = var.XResourceType
}
