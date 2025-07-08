provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_vpc_ports" "ports" {
  limit = var.limit
}
