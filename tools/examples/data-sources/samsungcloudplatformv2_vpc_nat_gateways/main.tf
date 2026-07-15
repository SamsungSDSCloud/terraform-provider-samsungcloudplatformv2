provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_vpc_nat_gateways" "natgateways" {
  size                   = var.size
  page                   = var.page
  sort                   = var.sort
  name                   = var.name
  vpc_id                 = var.vpc_id
  vpc_name               = var.vpc_name
  subnet_id              = var.subnet_id
  subnet_name            = var.subnet_name
  nat_gateway_ip_address = var.nat_gateway_ip_address
  state                  = var.state
}
