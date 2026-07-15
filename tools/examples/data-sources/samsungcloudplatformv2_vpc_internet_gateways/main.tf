provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_vpc_internet_gateways" "internetgateways" {
  size     = var.size
  page     = var.page
  sort     = var.sort
  id       = var.id
  name     = var.name
  type     = var.type
  state    = var.state
  vpc_id   = var.vpc_id
  vpc_name = var.vpc_name
}
