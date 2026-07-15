provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_organization_account" "account" {
  account_id = var.account_id
}