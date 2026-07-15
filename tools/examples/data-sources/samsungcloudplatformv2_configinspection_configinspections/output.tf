output "response" {
  value = {
    count : data.samsungcloudplatformv2_configinspection_configinspections.my_diagnosis_object_list.total_count,
    links : data.samsungcloudplatformv2_configinspection_configinspections.my_diagnosis_object_list.links,
    summary_responses : data.samsungcloudplatformv2_configinspection_configinspections.my_diagnosis_object_list.summary_responses
  }
}