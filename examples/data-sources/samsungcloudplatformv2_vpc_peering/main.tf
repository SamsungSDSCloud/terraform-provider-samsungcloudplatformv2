provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_vpc_peering" "vpc" {
        id = var.id
}
