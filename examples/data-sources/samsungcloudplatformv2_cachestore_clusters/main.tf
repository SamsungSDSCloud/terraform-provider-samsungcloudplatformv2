provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_cachestore_clusters" "clusters" {
  size = var.size
  page = var.page
  sort = var.sort
  name = var.name
  service_state = var.service_state
}
