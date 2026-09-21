provider "samsungcloudplatformv2" {
}

resource "samsungcloudplatformv2_ske_cluster_subnets" "subnets" {
  cluster_id               = var.cluster_id
  additional_subnet_id_list = var.additional_subnet_id_list
}
