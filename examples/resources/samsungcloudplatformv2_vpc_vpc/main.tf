provider "samsungcloudplatformv2" {
}

resource "samsungcloudplatformv2_vpc_vpc" "vpc" {
  cidr = var.vpc_cidr
  description = var.vpc_description
  name = var.vpc_name
}
