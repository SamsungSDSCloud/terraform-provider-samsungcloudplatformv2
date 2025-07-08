provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_vpc_nat_gateways" "natgateways" {
  limit = var.limit
}
