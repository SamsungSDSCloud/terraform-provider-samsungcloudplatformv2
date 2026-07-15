provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_iam_user_policy_bindings" "user_policy_bindings" {
  user_id = var.user_id
  size = var.size
  page = var.page
  sort = var.sort
  policy_id = var.policy_id
  policy_name = var.policy_name
  policy_type = var.policy_type
}
