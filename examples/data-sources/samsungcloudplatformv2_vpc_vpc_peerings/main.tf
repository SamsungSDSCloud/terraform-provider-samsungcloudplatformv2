provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_vpc_vpc_peerings" "vpcs" {
        id = var.id
        name = var.name
        requester_vpc_id = var.requester_vpc_id
        requester_vpc_name = var.requester_vpc_name
        approver_vpc_id = var.approver_vpc_id
        approver_vpc_name = var.approver_vpc_name
        account_type = var.account_type
        size = var.size
        sort = var.sort
        page = var.page
        state = var.state
}
