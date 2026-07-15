provider "samsungcloudplatformv2" {
}

resource "samsungcloudplatformv2_organization_account_remove" "remove" {
  organization_id    = var.organization_id
  account_id         = var.account_id
  target_account_ids = var.target_account_ids
}