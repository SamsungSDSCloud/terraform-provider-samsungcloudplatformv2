provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_scr_repositories" "repositories" {
  registry_id = var.registry_id
  sort        = "name:desc"
  size        = var.size
}
