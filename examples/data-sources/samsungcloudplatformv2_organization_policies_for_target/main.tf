provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_organization_policies_for_target" "policies" {
  target_id       = var.target_id
  organization_id = var.organization_id
  policy_category = var.policy_category
  name            = var.name
  size            = var.size
  page            = var.page
  sort            = var.sort
}
