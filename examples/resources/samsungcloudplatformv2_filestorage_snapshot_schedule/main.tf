provider "samsungcloudplatformv2" {
}

resource "samsungcloudplatformv2_filestorage_snapshot_schedule" "snapshot_schedule" {
  volume_id = var.volume_id
  snapshot_retention_count = var.snapshot_retention_count
  snapshot_schedule = var.snapshot_schedule
}