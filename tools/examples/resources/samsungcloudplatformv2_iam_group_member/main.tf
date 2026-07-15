provider "samsungcloudplatformv2" {
}

resource "samsungcloudplatformv2_iam_group_member" "group_member" {
  group_id = var.group_id
  user_id = var.user_id
}
