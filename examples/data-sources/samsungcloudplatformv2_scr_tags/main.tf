provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_scr_tags" "tags" {
  image_id = var.image_id
  page     = var.page
  size     = var.size
  sort     = var.sort
}
