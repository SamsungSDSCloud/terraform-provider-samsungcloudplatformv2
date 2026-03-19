provider "samsungcloudplatformv2" {
}

resource "samsungcloudplatformv2_loadbalancer_loadbalancer_private_nat_ip" "loadbalancerprivatenatip" {
  private_static_nat_create = var.static_nat
  loadbalancer_id = var.loadbalancer_id
}