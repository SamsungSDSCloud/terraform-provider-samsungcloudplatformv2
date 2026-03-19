provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_servicewatch_log_groups" "log_groups" {
  name = var.name
  retention_periods = [-1]
}
