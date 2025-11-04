provider "samsungcloudplatformv2" {
}

resource "samsungcloudplatformv2_vpc_transit_gateway_vpc_connection" "tgwvpcconnection" {
    transit_gateway_id = var.tgwvpcconnectiontransit_gateway_id
    vpc_id = var.tgwvpcconnection_vpc_id
}

