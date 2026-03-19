provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_servicewatch_log_group" "log_group" {
  log_group_id = var.id
}
