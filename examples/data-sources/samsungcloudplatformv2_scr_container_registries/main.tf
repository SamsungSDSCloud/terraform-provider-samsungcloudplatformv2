provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_scr_container_registries" "registries" {
  name = var.name
}
