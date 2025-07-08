provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_dns_records" "records" {
  hosted_zone_id = var.hosted_zone_id
}
