provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_vertica_cluster" "cluster" {
  id = var.id
}
