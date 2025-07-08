provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_vpc_internet_gateways" "internetgateways" {
  limit = var.limit
}
