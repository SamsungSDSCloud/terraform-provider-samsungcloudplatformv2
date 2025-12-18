provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_configinspection_configinspection" "my_diagnosis_object_detail" {
  diagnosis_id = var.diagnosis_id
}
