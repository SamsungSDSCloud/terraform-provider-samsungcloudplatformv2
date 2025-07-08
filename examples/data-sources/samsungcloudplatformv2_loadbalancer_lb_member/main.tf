provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_loadbalancer_lb_member" "lbmember" {
  id = var.id
  lb_server_group_id = var.lb_server_group_id
}