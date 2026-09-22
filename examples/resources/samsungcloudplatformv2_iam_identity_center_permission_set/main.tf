provider "samsungcloudplatformv2" {
}

resource "samsungcloudplatformv2_iam_identity_center_permission_set" "permission_set" {
  instance_id      = var.instance_id
  name             = var.name
  description      = var.description
  session_duration = var.session_duration
}
