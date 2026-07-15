provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_loadbalancer_lb_listener" "lb_listener" {
    id = var.id
}
