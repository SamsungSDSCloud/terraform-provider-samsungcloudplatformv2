provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_ske_cluster" "cluster" {
    id = var.id
}
