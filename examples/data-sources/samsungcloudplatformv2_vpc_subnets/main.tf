provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_vpc_subnets" "subnets" {
  limit = var.limit
}
