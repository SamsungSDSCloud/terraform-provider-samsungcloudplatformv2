provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_virtualserver_keypair" "keypair" {
  name = var.name

  filter {
    name = var.keypair_filter_name
    values = var.keypair_filter_values
    use_regex = var.keypair_filter_use_regex
  }
}
