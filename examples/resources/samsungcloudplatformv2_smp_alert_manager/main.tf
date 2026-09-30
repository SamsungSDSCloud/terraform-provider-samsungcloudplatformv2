provider "samsungcloudplatformv2" {
}

resource "samsungcloudplatformv2_smp_alert_manager" "alert_manager" {
  config_data           = var.config_data
  notification_group_id = var.notification_group_id
  workspace_id          = var.workspace_id
}
