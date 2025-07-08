provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_ske_nodepools" "nodepools" {
  cluster_id = var.cluster_id
}
