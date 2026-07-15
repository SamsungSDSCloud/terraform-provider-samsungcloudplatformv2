provider "samsungcloudplatformv2" {
}

resource "samsungcloudplatformv2_organization_service_control_policy" "scp_policy" {
  organization_id = var.organization_id
  name        = var.policy_name
  description = var.policy_description
  type        = var.policy_type
  document    = var.policy_document
}