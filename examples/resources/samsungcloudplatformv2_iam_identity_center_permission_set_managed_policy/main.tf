provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_iam_policy" "m" {
  id = var.managed_policy_id
}

resource "samsungcloudplatformv2_iam_identity_center_permission_set_managed_policy" "example" {
  permission_set_id   = var.permission_set_id
  instance_id         = var.instance_id
  managed_policy_id   = data.samsungcloudplatformv2_iam_policy.m.id
  managed_policy_name = data.samsungcloudplatformv2_iam_policy.m.policy.policy_name
}
