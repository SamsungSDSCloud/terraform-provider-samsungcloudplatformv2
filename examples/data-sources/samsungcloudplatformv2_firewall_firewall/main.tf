provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_firewall_firewall" "firewall" {
  id = var.id
  product_type = var.product_type
}
