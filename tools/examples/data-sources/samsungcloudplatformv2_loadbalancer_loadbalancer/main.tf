provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_loadbalancer_loadbalancer" "loadbalancer" {
    id = var.id
}
