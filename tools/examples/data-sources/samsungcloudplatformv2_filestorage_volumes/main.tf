provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_filestorage_volumes" "volumes" {
  limit = var.limit
  offset = var.offset
  sort = var.sort
  name = var.name
  type_name = var.type_name
}