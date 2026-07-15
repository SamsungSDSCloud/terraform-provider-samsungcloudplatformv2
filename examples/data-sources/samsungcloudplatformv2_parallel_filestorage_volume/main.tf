provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_parallel_filestorage_volume" "volume" {
  id = var.id
}