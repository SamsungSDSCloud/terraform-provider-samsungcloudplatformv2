provider "samsungcloudplatformv2" {
}

resource "samsungcloudplatformv2_loadbalancer_lb_health_check" "lbhealthcheck" {
  lb_health_check_create = var.lb_health_check_modify
}