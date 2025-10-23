provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_iam_policies" "policies" {
  size = var.size
  page = var.page
  sort = var.sort
  id = var.id
  policy_name = var.policy_name
}
