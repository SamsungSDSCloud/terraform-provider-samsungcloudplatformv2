provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_vertica_clusters" "clusters" {
  size = var.size
  page = var.page
  sort = var.sort
  name = var.name
  service_state = var.service_state
}
