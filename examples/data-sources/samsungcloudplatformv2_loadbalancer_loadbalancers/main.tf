provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_loadbalancer_loadbalancers" "loadbalancers" {
  name = var.name
  size = var.size
}
