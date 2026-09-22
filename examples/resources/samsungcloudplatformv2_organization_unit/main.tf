provider "samsungcloudplatformv2" {
}

resource "samsungcloudplatformv2_organization_unit" "org_unit" {
  name            = var.name
  description     = var.description
  organization_id = var.organization_id
  parent_unit_id  = var.parent_unit_id
  policy_ids      = var.policy_ids
}