provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_security_group_security_group_rule" "securitygrouprule" {
  id = var.id
}
