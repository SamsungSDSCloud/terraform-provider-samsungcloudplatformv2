provider "samsungcloudplatformv2" {
}

resource "samsungcloudplatformv2_organization_account_move" "move" {
  organization_id    = var.organization_id
  parent_unit_id     = var.parent_unit_id
  target_account_ids = var.target_account_ids
}