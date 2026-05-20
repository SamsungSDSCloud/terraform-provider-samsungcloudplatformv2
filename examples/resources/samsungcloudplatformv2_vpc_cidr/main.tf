provider "samsungcloudplatformv2" {
}

resource "samsungcloudplatformv2_vpc_cidr" "my_added_cidr_vpc" {
  vpc_id = var.vpc_id
  cidr   = var.cidr
}