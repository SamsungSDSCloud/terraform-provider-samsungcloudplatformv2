provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_smp_alert_manager" "alert_manager" {
  workspace_id = var.workspace_id
}
