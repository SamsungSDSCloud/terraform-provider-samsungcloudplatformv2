provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_multinodegpucluster_node_pools" "ids" {
  subnet_id         = var.subnet_id
  cluster_fabric_id = var.cluster_fabric_id
  node_pool_id      = var.node_pool_id
  zone              = var.zone
}
