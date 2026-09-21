provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_iam_identity_center_groups" "groups" {
  instance_id         = var.instance_id
  name                = var.name
  size                = var.size
  page                = var.page
  sort                = var.sort
  excluded_user_uuid  = var.excluded_user_uuid
  excluded_account_id = var.excluded_account_id
}
