provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_loadbalancer_lb_certificate" "lb_certificate" {
    id = var.id
}
