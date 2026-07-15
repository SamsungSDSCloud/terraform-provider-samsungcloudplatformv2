provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_sqlserver_cluster" "cluster" {
  id = var.id
}
