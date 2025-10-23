provider "samsungcloudplatformv2" {
}

resource "samsungcloudplatformv2_ske_nodepool" "nodepool" {
  name = var.name
  cluster_id = var.cluster_id
  image_os = var.image_os
  image_os_version = var.image_os_version
  is_auto_recovery = var.is_auto_recovery
  is_auto_scale = var.is_auto_scale
  min_node_count = var.min_node_count
  max_node_count = var.max_node_count
  keypair_name = var.keypair_name
  kubernetes_version = var.kubernetes_version
  desired_node_count = var.desired_node_count
  server_type_id = var.server_type_id
  volume_type_name = var.volume_type_name
  volume_size = var.volume_size
  labels = var.labels
  taints = var.taints
  server_group_id = var.server_group_id
  advanced_settings = var.advanced_settings
}
