provider "samsungcloudplatformv2" {
}

resource "samsungcloudplatformv2_loadbalancer_lb_member" "lbmember" {
  lb_member_create = var.lb_member
  lb_server_group_id = "22"
}
# resource "samsungcloudplatformv2_loadbalancer_lb_member" "lbmember" {
#   lb_member_create = var.lb_member_modify
#   lb_server_group_id = "22"
# }