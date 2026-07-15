provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_ske_nodepool" "nodepool" {
    id = var.id
}
