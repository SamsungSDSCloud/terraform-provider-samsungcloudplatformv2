provider "samsungcloudplatformv2" {
}

resource "samsungcloudplatformv2_organization" "org" {
  name                  = var.name
  use_scp_yn            = var.use_scp_yn
}