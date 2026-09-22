provider "samsungcloudplatformv2" {
}

resource "samsungcloudplatformv2_organization_delegation_account" "delegation_account" {
  account_id      = var.account_id
  organization_id = var.organization_id
  service_type    = var.service_type
}