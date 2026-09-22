provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_iam_identity_center_group" "example" {
  id          = var.id
  instance_id = var.instance_id
}