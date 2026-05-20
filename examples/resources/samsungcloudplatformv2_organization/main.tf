provider "samsungcloudplatformv2" {
}

resource "samsungcloudplatformv2_organization" "org" {
  name                  = var.name
  delegation_account_id = var.delegation_account_id
  use_scp_yn            = var.use_scp_yn
}