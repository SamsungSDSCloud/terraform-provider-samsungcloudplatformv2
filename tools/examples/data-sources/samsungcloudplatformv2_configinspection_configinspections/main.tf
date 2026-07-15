provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_configinspection_configinspections" "my_diagnosis_object_list" {
  with_count = true
  limit      = 5
}
