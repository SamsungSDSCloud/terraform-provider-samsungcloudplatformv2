provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_smp_alert_managers" "alert_managers" {
  workspace_id = var.workspace_id
}
