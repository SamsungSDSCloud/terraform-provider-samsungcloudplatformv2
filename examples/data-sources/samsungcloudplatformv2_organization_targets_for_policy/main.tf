provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_organization_targets_for_policy" "targets" {
  policy_id       = var.policy_id
  target_type     = var.target_type
  organization_id = var.organization_id
  policy_category = var.policy_category
  name            = var.name
  size            = var.size
  page            = var.page
  sort            = var.sort
}
