provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_organization" "org" {
  id = var.id
}
