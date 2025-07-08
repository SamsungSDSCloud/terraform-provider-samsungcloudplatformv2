provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_vpc_vpc_endpoints" "vpcendpoints" {
  limit = var.vpc_endpoint_limit
}
