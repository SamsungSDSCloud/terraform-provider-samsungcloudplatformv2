provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_vpc_private_nat_ips" "privatenatips" {
  private_nat_id = var.private_nat_id
  size = var.private_nat_ips_size
}
