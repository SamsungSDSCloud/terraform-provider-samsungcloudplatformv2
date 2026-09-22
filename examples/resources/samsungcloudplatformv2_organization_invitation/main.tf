provider "samsungcloudplatformv2" {
}

resource "samsungcloudplatformv2_organization_invitation" "invitation" {
  organization_id  = var.organization_id
  target_login_ids = var.target_login_ids
}