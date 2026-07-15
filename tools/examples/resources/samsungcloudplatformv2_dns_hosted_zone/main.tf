provider "samsungcloudplatformv2" {
}

resource "samsungcloudplatformv2_dns_hosted_zone" "hosted_zone" {
  hosted_zone_create = var.hosted_zone
  tags = var.tag
}