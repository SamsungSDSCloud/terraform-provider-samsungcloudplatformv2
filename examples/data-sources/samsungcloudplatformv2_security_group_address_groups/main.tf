provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_security_group_address_groups" "address_groups" {
  size = var.size
  page = var.page
  sort = var.sort
  id   = var.id
  name = var.name
}
