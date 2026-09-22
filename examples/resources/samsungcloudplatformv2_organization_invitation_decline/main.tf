provider "samsungcloudplatformv2" {
}

resource "samsungcloudplatformv2_organization_invitation_decline" "decline" {
  id       = var.invitation_id
}