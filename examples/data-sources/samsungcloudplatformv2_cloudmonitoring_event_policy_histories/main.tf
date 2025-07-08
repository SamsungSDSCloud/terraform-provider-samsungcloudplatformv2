provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_cloudmonitoring_event_policy_histories" "eventPolicyHistories" {
  query_end_dt = var.QueryEndDt
  query_start_dt = var.QueryStartDt
  x_resource_type = var.XResourceType
  event_policy_id = var.EventPolicyId

}
