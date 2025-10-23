provider "samsungcloudplatformv2" {
}

resource "samsungcloudplatformv2_iam_user" "user"{
  account_id = var.account_id
  description = var.description
  password = var.password
  temporary_password = var.temporary_password
  user_name = var.user_name
  tags = var.tags
  group_ids = var.group_ids
  policy_ids = var.policy_ids
  password_reuse_count = var.password_reuse_count
}
