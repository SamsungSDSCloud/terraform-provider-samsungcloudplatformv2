provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_iam_group_policy_bindings" "group_policy_bindings" {
  size = var.size
  page = var.page
  sort = var.sort
  group_id = var.group_id
  policy_id = var.policy_id
  policy_name = var.policy_name
  policy_type = var.policy_type
}
