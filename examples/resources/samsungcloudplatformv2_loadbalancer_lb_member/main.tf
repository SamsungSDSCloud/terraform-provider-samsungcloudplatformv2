provider "samsungcloudplatformv2" {
}

resource "samsungcloudplatformv2_loadbalancer_lb_member" "lbmember" {
  lb_member_create = var.lb_member_modify
  lb_server_group_id = "11"
}
# resource "samsungcloudplatformv2_loadbalancer_lb_member" "lbmember" {
#   lb_member_create = var.lb_member_modify  -> 수정 시 변경 (variables 참조)
#   lb_server_group_id = "22"
# }