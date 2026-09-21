provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_organization_unit" "org_unit" {
  unit_id = var.unit_id
  organization_id = var.organization_id
}