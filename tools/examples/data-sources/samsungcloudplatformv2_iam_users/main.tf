provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_iam_users" "users" {
  size = var.size
  page = var.page
  sort = var.sort
  email = var.email
  user_name = var.user_name
  type = var.type
  account_id = var.account_id
}
