provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_multinodegpucluster_cluster_fabric" "cluster_fabric" {
  id = var.id
}
