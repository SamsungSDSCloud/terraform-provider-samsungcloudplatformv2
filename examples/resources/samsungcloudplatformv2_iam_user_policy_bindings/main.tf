provider "samsungcloudplatformv2" {
}

resource "samsungcloudplatformv2_iam_user_policy_bindings" "user_policy_bindings"{
  user_id = var.user_id
  policy_ids = var.policy_ids
}
