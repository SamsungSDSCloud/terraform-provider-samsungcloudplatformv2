provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_servicewatch_alert" "alert" {
  id = var.id
}
