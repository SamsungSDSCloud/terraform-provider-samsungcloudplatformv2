provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_loadbalancer_lb_listeners" "lblisteners" {
  name = var.name
  size = var.size
}
