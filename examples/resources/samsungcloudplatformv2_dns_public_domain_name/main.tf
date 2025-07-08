provider "samsungcloudplatformv2" {
}

resource "samsungcloudplatformv2_dns_public_domain_name" "public_domain_name" {
  public_domain_name_create = var.public_domain_name
  tags = var.tag
}