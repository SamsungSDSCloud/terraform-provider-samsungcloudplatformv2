provider "samsungcloudplatformv2" {
}

resource "samsungcloudplatformv2_iam_group_policy_bindings" "group_policy_bindings" {
  group_id = var.group_id
  policy_ids = var.policy_ids
}
