provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_mysql_cluster" "cluster" {
  id = var.id
}
