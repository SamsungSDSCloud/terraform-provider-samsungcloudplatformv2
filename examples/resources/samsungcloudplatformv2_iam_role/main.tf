provider "samsungcloudplatformv2" {
}

resource "samsungcloudplatformv2_iam_role" "role"{
  name = var.name
  description = var.description
  max_session_duration = var.max_session_duration
  assume_role_policy_document = var.assume_role_policy_document
  principals = var.principals
  policy_ids = var.policy_ids
  tags = var.tags
}
