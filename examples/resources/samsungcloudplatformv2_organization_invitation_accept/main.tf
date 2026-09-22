provider "samsungcloudplatformv2" {
}

resource "samsungcloudplatformv2_organization_invitation_accept" "accept" {
  id       = var.invitation_id
}