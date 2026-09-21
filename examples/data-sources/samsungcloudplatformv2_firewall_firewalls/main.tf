provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_firewall_firewalls" "firewalls" {
  product_type = var.product_type
}
