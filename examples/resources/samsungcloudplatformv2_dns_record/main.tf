provider "samsungcloudplatformv2" {
}

resource "samsungcloudplatformv2_dns_record" "record" {
  hosted_zone_id = var.id
  record_create = var.record
}