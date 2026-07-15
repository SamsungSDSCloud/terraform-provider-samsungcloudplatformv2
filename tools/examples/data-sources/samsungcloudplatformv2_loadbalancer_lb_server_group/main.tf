provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_loadbalancer_lb_server_group" "lbservergroup" {
  id = var.id
}