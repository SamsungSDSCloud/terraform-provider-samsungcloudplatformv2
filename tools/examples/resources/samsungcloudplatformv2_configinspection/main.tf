provider "samsungcloudplatformv2" {
}

resource "samsungcloudplatformv2_configinspection" "my_config_inspection" {
  account_id             = var.account_id
  csp_type               = var.csp_type
  diagnosis_account_id   = var.diagnosis_account_id
  diagnosis_check_type   = var.diagnosis_check_type
  diagnosis_id           = var.diagnosis_id
  diagnosis_name         = var.diagnosis_name
  diagnosis_type         = var.diagnosis_type
  plan_type              = var.plan_type
  auth_key_request       = var.auth_key_request
  schedule_request       = var.schedule_request
  tags                   = var.tags
}