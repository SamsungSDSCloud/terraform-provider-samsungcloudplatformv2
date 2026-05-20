provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_vpc_private_nat" "my_private_nat" {
  private_nat_id = var.private_nat_id
}
