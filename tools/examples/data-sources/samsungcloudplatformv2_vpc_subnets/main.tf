provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_vpc_subnets" "subnets" {
  cidr     = var.cidr
  id       = var.id
  name     = var.name
  page     = 3
  size     = 4
  sort     = var.sort
  state    = var.state
  type     = ["LOCAL", "GENERAL", "VPC_ENDPOINT"]
  vpc_id   = var.vpc_id
  vpc_name = var.vpc_name
}
