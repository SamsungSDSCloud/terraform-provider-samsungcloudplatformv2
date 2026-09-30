provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_smp_workspace" "workspace" {
  workspace_id = var.workspace_id
}
