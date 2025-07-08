provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_filestorage_volume" "volume" {
  id = var.id
}