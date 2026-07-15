provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_iam_role_policy_bindings" "role_policy_bindings" {
  role_id = var.role_id
  size = var.size
  page = var.page
  sort = var.sort
  policy_id = var.policy_id
  policy_name = var.policy_name
  policy_type = var.policy_type
}
