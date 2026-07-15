provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_servicewatch_dashboards" "dashboards" {
  type = var.type
  favorite_enabled = var.favorite_enabled
}
