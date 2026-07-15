provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_mariadb_cluster" "cluster" {
  id = var.id
}
