provider "samsungcloudplatformv2" {
}

provider "samsungcloudplatformv2" {
  alias = "target_account"
  access_key = var.target_access_key
  secret_key = var.target_secret_key
}

resource "samsungcloudplatformv2_invitation_accept" "accept" {
  provider = samsungcloudplatformv2.target_account
  id       = var.invitation_id
}