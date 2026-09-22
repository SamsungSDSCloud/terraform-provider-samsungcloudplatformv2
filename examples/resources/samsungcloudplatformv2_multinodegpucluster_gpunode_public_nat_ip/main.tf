provider "samsungcloudplatformv2" {
  default_region = "kr-west1"
}

resource "samsungcloudplatformv2_multinodegpucluster_gpunode_public_nat_ip" "public_nat_ip" {
  gpu_node_id          = var.gpu_node_id
  public_ip_address_id = var.public_ip_address_id
  timeouts {
    create = var.create_timeouts
    delete = var.delete_timeouts
  }
}
