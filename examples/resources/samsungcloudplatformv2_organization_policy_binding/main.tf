provider "samsungcloudplatformv2" {
}

resource "samsungcloudplatformv2_organization_policy_binding" "policy_binding" {
  organization_id = var.organization_id
  entity          = var.entity
  target_ids      = var.target_ids
  policy_ids      = var.policy_ids
}
