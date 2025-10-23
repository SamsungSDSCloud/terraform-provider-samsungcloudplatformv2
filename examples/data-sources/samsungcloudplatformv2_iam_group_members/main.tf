provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_iam_group_members" "group_members" {
  size = var.size
  page = var.page
  sort = var.sort
  group_id = var.group_id
  user_name = var.user_name
  user_email = var.user_email
  creator_name = var.creator_name
  creator_email = var.creator_email
}
