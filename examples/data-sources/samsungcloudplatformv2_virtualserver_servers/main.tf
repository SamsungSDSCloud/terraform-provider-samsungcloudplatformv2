  provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_virtualserver_servers" "ids" {
  name = var.name
  ip = var.ip
  state = var.state
  product_category = var.product_category
  vpc_id = var.vpc_id
  server_type_id = var.server_type_id
  auto_scaling_group_id = var.auto_scaling_group_id

  filter {
    name = var.servers_filter_name
    values = var.servers_filter_values
    use_regex = var.servers_filter_use_regex
  }
}
