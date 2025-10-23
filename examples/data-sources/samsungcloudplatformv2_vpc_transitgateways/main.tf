provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_transitgateways" "vpctransitgateway" {
  size = var.size
  sort = var.sort
  page = var.page
  state = var.state
}
