provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_multinodegpucluster_gpunode_products" "products" {
  type     = var.type
  image_id = var.image_id
}
