provider "samsungcloudplatformv2" {
}

resource "samsungcloudplatformv2_multinodegpucluster_cluster_fabric_member" "cluster_fabric_member" {
  before_cluster_fabric_id = var.before_cluster_fabric_id
  after_cluster_fabric_id  = var.after_cluster_fabric_id
  gpu_node_id_list         = var.gpu_node_id_list
}
