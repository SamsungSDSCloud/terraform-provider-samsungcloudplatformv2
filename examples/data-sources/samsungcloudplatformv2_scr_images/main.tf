provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_scr_images" "images" {
  repository_id = var.repository_id
  name          = var.name
  sort          = var.sort
  page          = var.page
  size          = var.size
}
