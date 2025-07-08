provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_security_group_security_group_rules" "ids" {
  security_group_id = var.security_group_id
  size = var.size
  page = var.page
  id = var.id
  remote_ip_prefix = var.remote_ip_prefix
  remote_group_id = var.remote_group_id
  description = var.description
  direction = var.direction
  service = var.service
}
