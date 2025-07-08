provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_dns_record" "record" {
  hosted_zone_id = var.hosted_zone_id
  id = var.id
}
