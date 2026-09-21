provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_security_group_address_group_cidrs" "addresses" {
  address_group_id = var.address_group_id
  size             = var.size
  page             = var.page
  sort             = var.sort
  address          = var.address
}
