provider "samsungcloudplatformv2" {
}

resource "samsungcloudplatformv2_certificate_manager_self_sign" "certificatemanager01" {
  cn             = var.cn
  name                  = var.name
  not_after_dt           = var.not_after_dt
  not_before_dt                = var.not_before_dt
  organization                = var.organization
  region                = var.region
  tags                  = var.tags
  recipients            = var.recipients
  timezone              = var.timezone
}
