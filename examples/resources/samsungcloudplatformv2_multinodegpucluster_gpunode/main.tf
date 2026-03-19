provider "samsungcloudplatformv2" {
}

resource "samsungcloudplatformv2_multinodegpucluster_gpunode" "gpunode" {
  cluster_fabric_details = var.cluster_fabric_details
  gpu_node_name_prefix   = var.gpu_node_name_prefix
  image_id               = var.image_id
  os_user_password       = var.os_user_password
  region_id              = var.region_id
  server_details         = var.server_details
  server_type_id         = var.server_type_id
  subnet_id              = var.subnet_id
  vpc_id                 = var.vpc_id
  tags                   = var.tags
  timeouts {
    create = var.create_timeouts
    delete = var.delete_timeouts
  }
}