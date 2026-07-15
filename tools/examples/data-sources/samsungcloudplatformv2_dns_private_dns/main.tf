provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_dns_private_dns" "private_dns" {
  id = var.id
}
