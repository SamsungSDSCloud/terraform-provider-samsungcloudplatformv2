provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_iam_identity_center_permission_sets" "permission_sets" {
  instance_id = var.instance_id
  name        = var.name
  size        = var.size
  page        = var.page
  sort        = var.sort
}
