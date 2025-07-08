  provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_virtualserver_server" "server" {
  id = var.id

  name = var.name
  ip = var.ip
  state = var.state
  product_category = var.product_category
  vpc_id = var.vpc_id
  server_type_id = var.server_type_id
  auto_scaling_group_id = var.auto_scaling_group_id

  filter {
    name = var.server_filter_name
    values = var.server_filter_values
    use_regex = var.server_filter_use_regex
  }
}
