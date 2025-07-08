provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_filestorage_snapshot_schedule" "snapshot_schedule" {
  volume_id = var.volume_id
}