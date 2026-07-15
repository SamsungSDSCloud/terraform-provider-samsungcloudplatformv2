provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_firewall_firewalls" "ids" {
  product_type = var.product_type
}
