provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_organization_accounts" "accounts" {
  organization_id   = var.organization_id
  size              = var.size
  page              = var.page
  sort              = var.sort
  id                = var.id
  name              = var.name
  email             = var.email
  login_id          = var.login_id
  joined_start_date = var.joined_start_date
  joined_end_date   = var.joined_end_date
  joined_method     = var.joined_method
  exclude_policy_id = var.exclude_policy_id
  parent_unit_id    = var.parent_unit_id
  parent_unit_name  = var.parent_unit_name
  type              = var.type
}