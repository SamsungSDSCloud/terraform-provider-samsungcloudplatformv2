provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_scr_repository" "repository" {
  id = var.id
}
