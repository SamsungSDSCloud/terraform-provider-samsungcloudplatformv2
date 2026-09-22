provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_organization_service_control_policy" "scp_policy" {
  id       = var.policy_id
  organization_id = var.organization_id
}