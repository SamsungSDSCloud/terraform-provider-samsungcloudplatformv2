output "response" {
  value = {
    auth_key_responses : data.samsungcloudplatformv2_configinspection_configinspection.my_diagnosis_object_detail.auth_key_responses
    schedule_response : data.samsungcloudplatformv2_configinspection_configinspection.my_diagnosis_object_detail.schedule_response
    summary_responses : data.samsungcloudplatformv2_configinspection_configinspection.my_diagnosis_object_detail.summary_responses
  }
}