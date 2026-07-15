output "response" {
  value = {
    checklist_name       = data.samsungcloudplatformv2_configinspection_diagnosis.my_diagnosis_result.checklist_name
    count                = data.samsungcloudplatformv2_configinspection_diagnosis.my_diagnosis_result.total_count
    diagnosis_account_id = data.samsungcloudplatformv2_configinspection_diagnosis.my_diagnosis_result.diagnosis_account_id
    diagnosis_check_type = data.samsungcloudplatformv2_configinspection_diagnosis.my_diagnosis_result.diagnosis_check_type
    diagnosis_name       = data.samsungcloudplatformv2_configinspection_diagnosis.my_diagnosis_result.diagnosis_name
    links                = data.samsungcloudplatformv2_configinspection_diagnosis.my_diagnosis_result.links
    proceed_date         = data.samsungcloudplatformv2_configinspection_diagnosis.my_diagnosis_result.proceed_date
    result_detail_list   = data.samsungcloudplatformv2_configinspection_diagnosis.my_diagnosis_result.result_detail_list
    total                = data.samsungcloudplatformv2_configinspection_diagnosis.my_diagnosis_result.total
  }
}