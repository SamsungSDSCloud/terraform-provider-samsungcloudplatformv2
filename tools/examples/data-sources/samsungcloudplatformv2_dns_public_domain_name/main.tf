provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_dns_public_domain_name" "public_domain_name" {
  id = var.id
}
