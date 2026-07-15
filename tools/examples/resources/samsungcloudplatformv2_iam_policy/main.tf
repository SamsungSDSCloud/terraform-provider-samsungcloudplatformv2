provider "samsungcloudplatformv2" {
}

resource "samsungcloudplatformv2_iam_policy" "policy"{
  policy_name = var.policy_name
  description = var.policy_description
  tags = var.policy_tags
  policy_version = var.policy_version
}