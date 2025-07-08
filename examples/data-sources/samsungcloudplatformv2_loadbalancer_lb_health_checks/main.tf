provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_loadbalancer_lb_health_checks" "lbhealthchecks" {
  name = var.name
  size = var.size
}
