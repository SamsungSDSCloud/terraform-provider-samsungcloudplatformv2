provider "samsungcloudplatformv2" {
}

resource "samsungcloudplatformv2_ske_cluster_deletion_protection" "this" {
  cluster_id          = var.cluster_id
  deletion_protection_enabled = var.deletion_protection_enabled
}
