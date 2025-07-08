provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_virtualserver_server_groups" "ids" {
  filter {
    name = var.server_groups_filter_name
    values = var.server_groups_filter_values
    use_regex = var.server_groups_filter_use_regex
  }
}