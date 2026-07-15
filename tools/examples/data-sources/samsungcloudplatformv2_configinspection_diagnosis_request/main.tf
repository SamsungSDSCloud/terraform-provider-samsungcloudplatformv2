provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_configinspection_diagnosis_request" "my_diagnosis_request" {
  access_key           = var.access_key
  diagnosis_check_type = var.diagnosis_check_type
  diagnosis_id         = var.diagnosis_id
  secret_key           = var.secret_key
  tenant_id            = var.tenant_id
}
