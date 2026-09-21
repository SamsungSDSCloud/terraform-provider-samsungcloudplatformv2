provider "samsungcloudplatformv2" {
}

resource "samsungcloudplatformv2_iam_identity_center_permission_set_inline_policy" "example" {
  permission_set_id = var.permission_set_id
  instance_id       = var.instance_id
  policy_document   = jsonencode(var.policy_document)
}
