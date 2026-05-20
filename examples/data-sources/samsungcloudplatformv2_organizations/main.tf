provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_organizations" "orgs" {
  size              = var.size
  page              = var.page
  sort              = var.sort
  name              = var.name
  master_account_id = var.master_account_id
  organization_id   = var.organization_id
}