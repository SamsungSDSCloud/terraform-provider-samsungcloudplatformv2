provider "samsungcloudplatformv2" {
}

resource "samsungcloudplatformv2_ske_cluster" "cluster" {
  kubernetes_version = var.cluster_kubernetes_version
  name = var.cluster_name
  security_group_id_list = var.cluster_security_group_id_list
  default_subnet_id = var.cluster_default_subnet_id
  vpc_id = var.cluster_vpc_id
  nfs_volume_id = var.cluster_nfs_volume_id
  additional_subnet_id_list = var.cluster_additional_subnet_id_list
  deletion_protection_enabled = var.cluster_deletion_protection_enabled
  linked_resources = var.cluster_linked_resources
  private_endpoint_access_control_resources = var.private_endpoint_access_control_resources
  public_endpoint_access_control_ip = var.public_endpoint_access_control_ip
  service_watch_logging_enabled = var.service_watch_logging_enabled
  tags = {
    "terraform_key" = "terraform_value"
  }
}
