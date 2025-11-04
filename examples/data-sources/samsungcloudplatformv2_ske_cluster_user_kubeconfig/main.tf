provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_ske_cluster_user_kubeconfig" "cluster_user_kubeconfig" {
    cluster_id = var.cluster_id
    kubeconfig_type = var.kubeconfig_type
}
