provider "samsungcloudplatformv2" {
}

resource "samsungcloudplatformv2_loadbalancer_lb_server_group" "lbservergroup" {
  lb_server_group_create = var.lb_server_group
}