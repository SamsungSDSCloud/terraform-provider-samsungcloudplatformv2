provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_transitgateway" "vpctransitgateway" {
  id = var.id
}
