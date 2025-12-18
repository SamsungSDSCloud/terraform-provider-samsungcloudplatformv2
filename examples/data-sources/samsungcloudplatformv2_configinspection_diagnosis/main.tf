provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_configinspection_diagnosis" "my_diagnosis_result" {
  diagnosis_id               = var.diagnosis_id
  diagnosis_request_sequence = var.diagnosis_request_sequence
  with_count                 = var.with_count
  limit                      = var.limit
  marker                     = var.marker
  sort                       = var.sort
}
