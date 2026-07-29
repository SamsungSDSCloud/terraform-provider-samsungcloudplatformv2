provider "samsungcloudplatformv2" {
}

resource "samsungcloudplatformv2_organization_invitation_cancel" "cancel" {
  organization_id = var.organization_id
  ids             = var.ids
}