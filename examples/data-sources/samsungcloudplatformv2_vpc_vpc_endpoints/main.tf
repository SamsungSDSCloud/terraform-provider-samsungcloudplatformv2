provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_vpc_vpc_endpoints" "vpcendpoints" {
  id = var.id
  name = var.name
  vpc_id = var.vpc_id
  vpc_name = var.vpc_name
  resource_type = var.resource_type
  resource_key = var.resource_key
  endpoint_ip_address = var.endpoint_ip_address
  state = var.state
  limit = var.limit
  marker = var.marker
  sort = var.sort
}

