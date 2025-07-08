provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_ske_kubernetes_versions" "kubernetes_versions" {
    region = var.region
}
