provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_ske_nodepoolnodes" "nodepoolnodes" {
  nodepool_id = var.nodepool_id
}