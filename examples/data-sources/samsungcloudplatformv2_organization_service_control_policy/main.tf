provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_organization_service_control_policy" "scp_policy" {
  policy_id       = var.policy_id
  organization_id = var.organization_id
}