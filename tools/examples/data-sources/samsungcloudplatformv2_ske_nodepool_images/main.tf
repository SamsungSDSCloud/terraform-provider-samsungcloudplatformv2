provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_ske_nodepool_images" "nodepool_images" {
  scp_original_image_type = var.scp_original_image_type
  kubernetes_version = var.kubernetes_version
  os = var.os
  size = var.size
  page = var.page
  sort = var.sort
}
