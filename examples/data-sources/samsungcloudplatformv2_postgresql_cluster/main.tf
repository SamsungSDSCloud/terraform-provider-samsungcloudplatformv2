provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_postgresql_cluster" "cluster" {
  id = var.id
}
