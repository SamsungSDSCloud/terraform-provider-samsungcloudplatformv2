provider "samsungcloudplatformv2" {
}

resource "samsungcloudplatformv2_scr_container_registry" "registry" {
  name                    = var.name
  public_visible_enabled  = var.public_visible_enabled
  public_endpoint_enabled = var.public_endpoint_enabled
  private_acl_enabled     = var.private_acl_enabled
  public_acl_enabled      = var.public_acl_enabled

  private_acl_resources = var.private_acl_resources

  public_acl_resources = var.public_acl_resources
}
