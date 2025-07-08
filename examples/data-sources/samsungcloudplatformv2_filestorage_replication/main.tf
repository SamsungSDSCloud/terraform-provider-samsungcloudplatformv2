provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_filestorage_replication" "replication" {
    id = var.id
    volume_id = var.volume_id
}