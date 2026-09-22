provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_organization_delegation_accounts" "delegation_accounts" {
  organization_id = var.organization_id
  service_type    = var.service_type
  account_id      = var.account_id
  size            = var.size
  page            = var.page
  sort            = var.sort
}