provider "samsungcloudplatformv2" {
}

resource "samsungcloudplatformv2_iam_identity_center_user" "example" {
  instance_id = var.instance_id
  user_id     = var.user_id
  name        = var.name
  email       = var.email
  password    = var.password
}