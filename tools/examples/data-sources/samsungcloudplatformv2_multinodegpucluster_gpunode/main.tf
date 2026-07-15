provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_multinodegpucluster_gpunode" "gpunode" {
  id = var.id
}