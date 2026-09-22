provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_multinodegpucluster_cluster_fabrics" "ids" {
  cluster_fabric_name = var.cluster_fabric_name
  state               = var.state
  node_pool_id        = var.node_pool_id
}
