provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_epas_cluster" "cluster" {
  id = var.id
}
