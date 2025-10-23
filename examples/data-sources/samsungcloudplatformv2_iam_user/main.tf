provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_iam_user" "user" {
  account_id = var.account_id
  user_id = var.user_id
}
