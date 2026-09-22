provider "samsungcloudplatformv2" {
}

resource "samsungcloudplatformv2_security_group_address_group" "address_group" {
  name        = var.name
  description = var.description
  addresses   = var.addresses
  tags        = var.tags
}
