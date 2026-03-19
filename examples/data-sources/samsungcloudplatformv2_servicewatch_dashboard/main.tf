provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_servicewatch_dashboard" "dashboard" {
  id = var.id
}
