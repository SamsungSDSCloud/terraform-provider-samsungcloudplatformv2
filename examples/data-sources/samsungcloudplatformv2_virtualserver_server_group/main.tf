provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_virtualserver_server_group" "server_group" {
  id = var.id

  filter {
    name = var.server_group_filter_name
    values = var.server_group_filter_values
    use_regex = var.server_group_filter_use_regex
  }
}
