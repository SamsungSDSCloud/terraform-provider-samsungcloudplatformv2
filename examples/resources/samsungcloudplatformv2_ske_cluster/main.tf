provider "samsungcloudplatformv2" {
}

resource "samsungcloudplatformv2_ske_cluster" "cluster" {
  cloud_logging_enabled = var.cluster_cloud_logging_enabled
  kubernetes_version = var.cluster_kubernetes_version
  name = var.cluster_name
  security_group_id_list = var.cluster_security_group_id_list
  subnet_id = var.cluster_subnet_id
  vpc_id = var.cluster_vpc_id
  volume_id = var.cluster_volume_id
  private_endpoint_access_control_resources = var.private_endpoint_access_control_resources
  public_endpoint_access_control_ip = var.public_endpoint_access_control_ip
  service_watch_logging_enabled = var.service_watch_logging_enabled
  tags = {
      "terraform_key" = "terraform_value"
  }
}
