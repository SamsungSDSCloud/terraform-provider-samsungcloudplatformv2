provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_smp_workspace_configurations" "workspace_configurations" {
  workspace_id = var.workspace_id
}
