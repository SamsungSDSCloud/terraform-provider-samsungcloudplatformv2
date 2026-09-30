provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_smp_workspace_configuration" "workspace_configuration" {
  workspace_id = var.workspace_id
}
