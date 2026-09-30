provider "samsungcloudplatformv2" {
}

resource "samsungcloudplatformv2_smp_workspace" "workspace" {
  name = var.workspace_name

  private_acl_enabled = true

  private_acl_resources = [
    {
      resource_id   = var.resource_id
      resource_name = var.resource_name
      resource_type = var.resource_type
    },
  ]
}
