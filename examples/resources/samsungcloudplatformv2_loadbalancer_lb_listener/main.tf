provider "samsungcloudplatformv2" {
}


# locals {
#   lb_listeners = {
#     "listener1" = { lb_listener_create = var.lb_listener1 },
#     "listener2" = { lb_listener_create = var.lb_listener2 },
#     "listener3" = { lb_listener_create = var.lb_listener3 },
#     "listener4" = { lb_listener_create = var.lb_listener4 },
#   }
# }
#
# resource "samsungcloudplatformv2_loadbalancer_lb_listener" "lblistener" {
#   for_each           = local.lb_listeners
#   lb_listener_create = each.value.lb_listener_create
# }


resource "samsungcloudplatformv2_loadbalancer_lb_listener" "lblistener" {
  lb_listener_create = var.lb_listener_tcp
}
