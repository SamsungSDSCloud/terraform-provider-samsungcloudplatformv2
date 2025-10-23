provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_baremetal_baremetals" "ids" {
  server_name = var.server_name
  state = var.state
  policy_ip = var.ip
  vpc_id = var.vpc_id

  filter {
    name = var.baremetals_filter_name
    values = var.baremetals_filter_values
    use_regex = var.baremetals_filter_use_regex
  }
}