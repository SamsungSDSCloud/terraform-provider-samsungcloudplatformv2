provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_loadbalancer_loadbalancer_private_nat_ip" "loadbalancerprivatenatip" {
    loadbalancer_id = var.loadbalancer_id
}