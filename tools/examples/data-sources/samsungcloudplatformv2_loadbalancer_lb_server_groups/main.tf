provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_loadbalancer_lb_server_groups" "lbservergroups" {
  name = var.name
  size = var.size
}
