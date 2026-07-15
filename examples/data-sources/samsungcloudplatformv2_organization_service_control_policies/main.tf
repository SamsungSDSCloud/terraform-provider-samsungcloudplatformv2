provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_organization_service_control_policies" "scp_policies" {
  organization_id   = var.organization_id
  name              = var.name
  type              = var.type
  id                = var.id
  size              = var.size
  page              = var.page
  sort              = var.sort
  exclude_target_id = var.exclude_target_id
}