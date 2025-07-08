provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_virtualserver_volumes" "ids" {
  name = var.name
  filter {
    name = var.volumes_filter_name
    values = var.volumes_filter_values
    use_regex = var.volumes_filter_use_regex
  }
}
