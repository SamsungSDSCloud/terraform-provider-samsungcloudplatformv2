provider "samsungcloudplatformv2" {
}

resource "samsungcloudplatformv2_servicewatch_log_stream" "log_stream" {
  log_group_id= var.log_group_id
  name = var.name
}