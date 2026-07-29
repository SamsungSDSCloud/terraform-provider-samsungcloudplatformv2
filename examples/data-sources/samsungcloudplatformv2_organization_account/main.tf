provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_organization_account" "account" {
  account_id      = var.account_id
  organization_id = var.organization_id != "" ? var.organization_id : null
  lazy_policy     = var.lazy_policy
}