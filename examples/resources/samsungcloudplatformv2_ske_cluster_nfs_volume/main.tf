provider "samsungcloudplatformv2" {
}

resource "samsungcloudplatformv2_ske_cluster_nfs_volume" "nfs_volume" {
  cluster_id    = var.cluster_id
  nfs_volume_id = var.nfs_volume_id
}
