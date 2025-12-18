output "response" {
  value = {
    count : data.samsungcloudplatformv2_configinspection_diagnoses.my_diagnosis_result_list.total_count,
    links : data.samsungcloudplatformv2_configinspection_diagnoses.my_diagnosis_result_list.links,
    diagnosis_result_responses : data.samsungcloudplatformv2_configinspection_diagnoses.my_diagnosis_result_list.diagnosis_result_responses
  }
}