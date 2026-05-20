provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_vpc_private_nats" "privatenats" {
  size                  = var.size
  page                  = var.page
  sort                  = var.sort
  name                  = var.name
  cidr                  = var.cidr
  vpc_id                = var.vpc_id
  service_resource_id   = var.service_resource_id
  service_type          = var.service_type
  service_resource_name = var.service_resource_name
  state                 = var.state
}
