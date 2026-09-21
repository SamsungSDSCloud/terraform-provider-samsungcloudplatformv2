provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_cloudcontrol_landing_zone" "landing_zone" {
  landing_zone_id = var.landing_zone_id
}