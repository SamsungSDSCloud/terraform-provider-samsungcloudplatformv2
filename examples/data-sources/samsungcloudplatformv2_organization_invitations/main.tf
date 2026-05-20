provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_organization_invitations" "invitations" {
  organization_id = var.organization_id
  size            = var.size
  page            = var.page
  sort            = var.sort
  account_id      = var.account_id
  account_name    = var.account_name
  account_email   = var.account_email
  state           = var.state
  login_id        = var.login_id
}