provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_loadbalancer_lb_members" "lbmembers" {
  lb_server_group_id = var.lb_server_group_id
  size = var.size
  member_ip = var.member_ip
  member_port = var.member_port
}
