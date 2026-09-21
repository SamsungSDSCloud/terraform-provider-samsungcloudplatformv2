provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_security_group_address_group" "address_group" {
  id = var.id
}
