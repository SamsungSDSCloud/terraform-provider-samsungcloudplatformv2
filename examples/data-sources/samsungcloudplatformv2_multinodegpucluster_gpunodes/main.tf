provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_multinodegpucluster_gpunodes" "ids" {
  gpu_node_name        = var.gpu_node_name
  state                = var.state
  ip                   = var.ip
  vpc_id               = var.vpc_id
  cluster_fabric_name  = var.cluster_fabric_name
  cluster_fabric_id    = var.cluster_fabric_id
}