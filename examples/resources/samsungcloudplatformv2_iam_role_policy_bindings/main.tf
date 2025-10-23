provider "samsungcloudplatformv2" {
}

resource "samsungcloudplatformv2_iam_role_policy_bindings" "role_policy_bindings"{
  role_id = var.role_id
  policy_ids = var.policy_ids
}
