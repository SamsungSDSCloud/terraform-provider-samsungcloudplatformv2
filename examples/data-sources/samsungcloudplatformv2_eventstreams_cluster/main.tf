provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_eventstreams_cluster" "cluster" {
  id = var.id
}
