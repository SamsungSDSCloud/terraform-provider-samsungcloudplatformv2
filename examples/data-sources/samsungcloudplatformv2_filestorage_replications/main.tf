provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_filestorage_replications" "replications" {
    volume_id = var.volume_id
}