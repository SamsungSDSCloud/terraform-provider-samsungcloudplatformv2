provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_scr_container_registry" "registry" {
  id = var.id
}
