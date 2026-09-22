provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_iam_identity_center_permission_set" "permission_set" {
  instance_id = var.instance_id
  id          = var.id
}