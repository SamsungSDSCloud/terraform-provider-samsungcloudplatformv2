provider "samsungcloudplatformv2" {
}

resource "samsungcloudplatformv2_certificate_manager" "certificatemanager01" {
  cert_body             = var.cert_body
  cert_chain            = var.cert_chain
  name                  = var.name
  private_key           = var.private_key
  region                = var.region
  tags                  = var.tags
  recipients            = var.recipients
  timezone              = var.timezone
}
