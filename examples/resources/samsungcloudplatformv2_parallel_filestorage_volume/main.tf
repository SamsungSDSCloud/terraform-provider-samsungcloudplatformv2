provider "samsungcloudplatformv2" {
}

resource "samsungcloudplatformv2_parallel_filestorage_volume" "volume" {
  name = var.name
  capacity_tb = var.capacity_tb
  tags = var.tags
  zone = var.zone
  access_rules = var.access_rules
}