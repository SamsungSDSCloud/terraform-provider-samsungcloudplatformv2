provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_scr_tag_packages" "packages" {
  tags_id = var.tags_id
  size    = var.size
}
