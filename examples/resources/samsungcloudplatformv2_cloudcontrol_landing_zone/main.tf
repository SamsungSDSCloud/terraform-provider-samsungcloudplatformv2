provider "samsungcloudplatformv2" {
}

resource "samsungcloudplatformv2_cloudcontrol_landing_zone" "landing_zone" {
  additional_ou_name          = var.additional_ou_name
  agree_yn                    = var.agree_yn
  audit_account_name          = var.audit_account_name
  audit_login_id              = var.audit_login_id
  basic_ou_name               = var.basic_ou_name
  detective_guardrail_status  = var.detective_guardrail_status
  log_archive_account_name    = var.log_archive_account_name
  log_archive_login_id        = var.log_archive_login_id
  sso_type                    = var.sso_type
}