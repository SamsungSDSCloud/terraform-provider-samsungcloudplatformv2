provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_servicewatch_log_stream" "log_stream" {
  log_group_id = var.log_group_id
  log_stream_id = var.log_stream_id
}
