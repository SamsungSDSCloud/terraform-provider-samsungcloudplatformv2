provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_organization_units" "org_units" {
  parent_unit_id    = var.parent_unit_id
  organization_id   = var.organization_id
  name              = var.name
  exclude_policy_id = var.exclude_policy_id
}