provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_iam_roles" "roles" {
  size = var.size
  page = var.page
  sort = var.sort
  name = var.name
  role_type = var.role_type
  account_id = var.account_id
}
