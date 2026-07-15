provider "samsungcloudplatformv2" {
}

resource "samsungcloudplatformv2_loadbalancer_loadbalancer_public_nat_ip" "loadbalancerpublicnatip" {
  static_nat_create = var.static_nat
  loadbalancer_id = var.loadbalancer_id
}