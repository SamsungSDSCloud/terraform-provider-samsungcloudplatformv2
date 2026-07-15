provider "samsungcloudplatformv2" {
}

resource "samsungcloudplatformv2_vpc_vpc_peering" "vpc_peering01" {
  approver_vpc_account_id = var.approver_vpc_account_id
  approver_vpc_id         = var.approver_vpc_id
  name                    = var.name
  requester_vpc_id        = var.requester_vpc_id
  tags                    = var.tags
  description             = var.description
}
