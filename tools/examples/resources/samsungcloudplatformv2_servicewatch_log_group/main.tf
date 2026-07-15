provider "samsungcloudplatformv2" {
}

resource "samsungcloudplatformv2_servicewatch_log_group" "log_group" {
  name = var.log_group_name
  retention_period = var.retention_period
}