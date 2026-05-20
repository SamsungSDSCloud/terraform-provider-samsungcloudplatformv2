provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_vpc_publicips" "publicips" {
  size                   = var.size
  page                   = var.page
  sort                   = var.sort
  ip_address             = var.ip_address
  state                  = var.state
  attached_resource_type = var.attached_resource_type
  attached_resource_id   = var.attached_resource_id
  attached_resource_name = var.attached_resource_name
  vpc_id                 = var.vpc_id
  type                   = var.type
}
