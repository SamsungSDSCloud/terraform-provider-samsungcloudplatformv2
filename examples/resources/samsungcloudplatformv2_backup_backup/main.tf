provider "samsungcloudplatformv2" {
}

resource "samsungcloudplatformv2_backup_backup" "backup" {
  region = var.region

  name               = var.name
  policy_category    = var.policy_category
  policy_type        = var.policy_type
  server_uuid        = var.server_uuid
  server_category    = var.server_category
  encrypt_enabled    = var.encrypt_enabled
  retention_period   = var.retention_period
  schedules          = var.schedules
  tags               = var.tags
}