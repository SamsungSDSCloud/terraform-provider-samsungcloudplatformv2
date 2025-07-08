provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_vpc_vpcs" "vpcs" {
  limit = var.limit
}
