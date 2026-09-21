provider "samsungcloudplatformv2" {
}

resource "samsungcloudplatformv2_iam_identity_center_permission_set_custom_policy" "example" {
  permission_set_id = var.permission_set_id
  instance_id       = var.instance_id
  name              = var.name
}
