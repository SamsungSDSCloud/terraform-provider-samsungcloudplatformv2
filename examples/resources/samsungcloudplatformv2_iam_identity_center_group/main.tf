provider "samsungcloudplatformv2" {
}

resource "samsungcloudplatformv2_iam_identity_center_group" "example" {
  instance_id = var.instance_id
  name        = var.name
  description = var.description
}