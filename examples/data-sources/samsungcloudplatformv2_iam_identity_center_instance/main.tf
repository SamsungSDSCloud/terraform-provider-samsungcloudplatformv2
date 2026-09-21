provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_iam_identity_center_instance" "instance" {
  id   = var.id
}