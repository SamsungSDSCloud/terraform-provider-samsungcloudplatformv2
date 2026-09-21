provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_scr_image" "image" {
  id = var.id
}
