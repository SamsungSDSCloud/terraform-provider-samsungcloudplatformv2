provider "samsungcloudplatformv2" {
}

resource "samsungcloudplatformv2_iam_group" "group" {
  name = var.group_name
  description = var.group_description
  tags = var.group_tags
  policy_ids = var.group_policy_ids
  user_ids = var.group_user_ids
}
