provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_vpc_transit_gateway_rules" "vpc_transitgateway_routingrules" {
    transit_gateway_id = var.transit_gateway_id
    size = var.size
    page = var.page
    sort = var.sort
    id = var.id
    tgw_connection_vpc_id = var.tgw_connection_vpc_id
    tgw_connection_vpc_name = var.tgw_connection_vpc_name
    source_type = var.source_type
    destination_type = var.destination_type
    destination_cidr = var.destination_cidr
    state = var.state
}
