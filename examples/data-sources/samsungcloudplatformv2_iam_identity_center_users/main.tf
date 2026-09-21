provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_iam_identity_center_users" "users" {
  instance_id          = var.instance_id
  user_id              = var.user_id
  size                 = var.size
  page                 = var.page
  sort                 = var.sort
  excluded_group_id    = var.excluded_group_id
  excluded_account_id  = var.excluded_account_id
}
