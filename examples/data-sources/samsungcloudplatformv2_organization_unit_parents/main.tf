provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_organization_unit_parents" "org_unit_parents" {
  unit_id = var.unit_id
}