provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_ske_clusters" "ids" {
    page               = var.page
    size               = var.size
    sort               = var.sort
    name               = var.name
    subnet_id          = var.subnet_id
    status             = var.status
    kubernetes_version = var.kubernetes_version
    region = var.clusters_region
}
