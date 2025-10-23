provider "samsungcloudplatformv2" {
}

resource "samsungcloudplatformv2_vpc_transitgateway_vpcconnection" "tgwvpcconnection" {
    transit_gateway_id = var.tgwvpcconnectiontransit_gateway_id
    vpc_id = var.tgwvpcconnection_vpc_id
}

