provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_iam_identity_center_group_member" "example" {
  instance_id = var.instance_id
  group_id    = var.group_id
}
