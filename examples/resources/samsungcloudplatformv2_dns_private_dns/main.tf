provider "samsungcloudplatformv2" {
}

resource "samsungcloudplatformv2_dns_private_dns" "private_dns" {
  private_dns_create = var.private_dns
  tags = var.tag
}