provider "samsungcloudplatformv2" {
}

resource "samsungcloudplatformv2_baremetal_blockstorage_volume" "volume"{
  region = var.region
  name = var.block_storage_name
  disk_type = var.disk_type
  size_gb = var.size_gb
  attachments = var.attachments
  qos = var.qos
  tags = var.tags
  timeouts {
    create = var.create_timeouts
    delete = var.delete_timeouts
  }
}
