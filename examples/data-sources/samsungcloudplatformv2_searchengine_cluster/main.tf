provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_searchengine_cluster" "cluster" {
  id = var.id
}
