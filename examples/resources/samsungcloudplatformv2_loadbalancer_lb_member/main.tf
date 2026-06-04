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

# resource "samsungcloudplatformv2_loadbalancer_lb_member" "lbmember" {
#   lb_server_group_id = "68d9e7ab66604f5190f85f8aefb4d037"
# 
#   lb_member_create = {
#     name         = "manual-ip-member"
#     object_type  = "MANUAL"
#     member_ip    = "10.0.1.50"
#     member_port  = 8080
#     member_weight = 1
#     member_state = "ENABLE"
#   }
# }