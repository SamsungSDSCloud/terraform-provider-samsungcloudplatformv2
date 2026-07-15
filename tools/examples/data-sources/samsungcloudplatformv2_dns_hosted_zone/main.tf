provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_dns_hosted_zone" "hosted_zone" {
  id = var.id
}
