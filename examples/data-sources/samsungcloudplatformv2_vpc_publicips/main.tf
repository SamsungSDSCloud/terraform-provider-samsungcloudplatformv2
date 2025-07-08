provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_vpc_publicips" "publicips" {
  limit = var.limit
}
