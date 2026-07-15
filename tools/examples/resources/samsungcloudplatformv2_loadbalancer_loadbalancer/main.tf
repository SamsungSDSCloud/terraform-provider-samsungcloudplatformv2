provider "samsungcloudplatformv2" {
}

resource "samsungcloudplatformv2_loadbalancer_loadbalancer" "loadbalancer" {
  loadbalancer_create = var.loadbalancer
}