provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_iam_identity_center_account_assignment" "example" {
  instance_id       = var.instance_id
  target_account_id = var.target_account_id
  principal_id      = var.principal_id
}