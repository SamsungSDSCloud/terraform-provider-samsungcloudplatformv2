provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_ske_cluster_deletion_protection" "this" {
  cluster_id = var.cluster_id
}
