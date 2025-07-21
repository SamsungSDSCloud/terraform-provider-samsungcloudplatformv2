provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_vpc_private_nats" "privatenats" {
  size = var.private_nats_size
}
