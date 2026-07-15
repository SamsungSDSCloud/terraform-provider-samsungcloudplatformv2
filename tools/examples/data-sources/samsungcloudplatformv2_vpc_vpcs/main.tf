provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_vpc_vpcs" "vpcs" {
  cidr  = var.cidr
  id    = var.id
  name  = var.name
  page  = var.page
  size  = var.size
  sort  = var.sort
  state = var.state
}
