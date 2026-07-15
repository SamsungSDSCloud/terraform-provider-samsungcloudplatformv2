provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_cachestore_cluster" "cluster" {
  id = var.id
}
