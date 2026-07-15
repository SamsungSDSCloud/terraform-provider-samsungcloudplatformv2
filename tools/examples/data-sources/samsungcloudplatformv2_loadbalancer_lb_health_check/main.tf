provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_loadbalancer_lb_health_check" "lbhealthcheck" {
  id = var.id
}