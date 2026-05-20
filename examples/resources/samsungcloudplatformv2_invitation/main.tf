provider "samsungcloudplatformv2" {
}

resource "samsungcloudplatformv2_invitation" "invitation" {
  organization_id  = var.organization_id
  target_login_ids = var.target_login_ids
}