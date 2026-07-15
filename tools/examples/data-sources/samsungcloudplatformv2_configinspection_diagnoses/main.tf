provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_configinspection_diagnoses" "my_diagnosis_result_list" {
  with_count = true
  limit      = 5
}