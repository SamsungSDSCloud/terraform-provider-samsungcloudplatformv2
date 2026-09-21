provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_iam_identity_center_permission_set_policies" "example" {
  instance_id       = var.instance_id
  permission_set_id = var.permission_set_id
}
