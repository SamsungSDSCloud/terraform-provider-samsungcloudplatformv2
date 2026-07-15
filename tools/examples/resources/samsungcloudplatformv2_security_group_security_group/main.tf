provider "samsungcloudplatformv2" {
}

resource "samsungcloudplatformv2_security_group_security_group" "securitygroup" {
  name = var.name
  description = var.description
  loggable = var.loggable
  tags = var.security_group_tags
}
