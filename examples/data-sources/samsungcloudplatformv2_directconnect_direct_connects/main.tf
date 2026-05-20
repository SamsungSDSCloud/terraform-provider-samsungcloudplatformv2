provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_directconnect_direct_connects" "directconnects" {
  size     = var.size
  page     = var.page
  sort     = var.sort
  id       = var.id
  name     = var.name
  state    = var.state
  vpc_id   = var.vpc_id
  vpc_name = var.vpc_name
}
