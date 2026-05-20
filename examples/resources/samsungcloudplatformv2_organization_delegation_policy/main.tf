provider "samsungcloudplatformv2" {
}

resource "samsungcloudplatformv2_organization_delegation_policy" "delegation_policy" {
  organization_id = var.organization_id
  document        = var.document
}