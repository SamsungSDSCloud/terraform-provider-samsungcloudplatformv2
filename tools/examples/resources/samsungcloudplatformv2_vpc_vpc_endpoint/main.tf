provider "samsungcloudplatformv2" {
}

resource "samsungcloudplatformv2_vpc_vpc_endpoint" "vpcendpoint" {
  name = var.vpcendpoint_name
  vpc_id = var.vpcendpoint_vpc_id
  subnet_id = var.vpcendpoint_subnet_id
  resource_type = var.vpcendpoint_resource_type
  resource_key = var.vpcendpoint_resource_resource_key
  resource_info = var.vpcendpoint_resource_resource_info
  endpoint_ip_address = var.vpcendpoint_endpoint_ip_address
  description = var.vpcendpoint_description
  tags = var.vpcendpoint_tags
}
