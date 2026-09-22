provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_filestorage_access_rules" "all" {
  file_storage_id = var.file_storage_id
}
