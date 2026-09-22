provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_multinodegpucluster_gpunode_images" "images" {
  region_id = var.region_id
}
