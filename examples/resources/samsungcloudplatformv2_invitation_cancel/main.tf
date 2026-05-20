provider "samsungcloudplatformv2" {
}

resource "samsungcloudplatformv2_invitation_cancel" "cancel" {
  organization_id = var.organization_id
  ids             = var.ids
}