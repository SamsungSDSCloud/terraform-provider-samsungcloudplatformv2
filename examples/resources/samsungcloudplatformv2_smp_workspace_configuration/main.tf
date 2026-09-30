provider "samsungcloudplatformv2" {
}

resource "samsungcloudplatformv2_smp_workspace_configuration" "workspace_configuration" {
  workspace_id     = var.workspace_id
  retention_period = var.retention_period
}
