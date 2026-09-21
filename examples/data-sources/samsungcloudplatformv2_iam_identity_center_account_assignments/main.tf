provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_iam_identity_center_account_assignments" "account_assignments" {
  instance_id          = var.instance_id
  target_account_id    = var.target_account_id
  size                 = var.size
  page                 = var.page
  sort                 = var.sort
  target_account_name  = var.target_account_name
  target_account_email = var.target_account_email
  permission_set_id    = var.permission_set_id
  principal_name       = var.principal_name
  role_srn             = var.role_srn
}
