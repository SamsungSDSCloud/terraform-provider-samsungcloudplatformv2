provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_virtualserver_keypairs" "names" {
  filter {
    name = var.keypairs_filter_name
    values = var.keypairs_filter_values
    use_regex = var.keypairs_filter_use_regex
  }
}
