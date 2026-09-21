provider "samsungcloudplatformv2" {
}

resource "samsungcloudplatformv2_iam_identity_center_account_assignment" "example" {
  instance_id       = var.instance_id
  target_account_id = var.target_account_id
  principal_id      = var.principal_id
  principal_type    = var.principal_type
  permission_set_id = var.permission_set_id
}
