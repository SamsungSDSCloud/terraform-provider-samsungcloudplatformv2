provider "samsungcloudplatformv2" {
}

resource "samsungcloudplatformv2_iam_identity_center_instance" "instance" {
  name                  = var.name
  description           = var.description
  identity_store_type   = var.identity_store_type
  identity_store_config = var.identity_store_config
  self_managed_password = var.self_managed_password
  region                = var.region
}