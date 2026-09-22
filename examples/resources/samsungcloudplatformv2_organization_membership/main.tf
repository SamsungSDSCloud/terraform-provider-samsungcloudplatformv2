provider "samsungcloudplatformv2" {
}

resource "samsungcloudplatformv2_organization_membership" "leave" {
  organization_id = var.organization_id
}