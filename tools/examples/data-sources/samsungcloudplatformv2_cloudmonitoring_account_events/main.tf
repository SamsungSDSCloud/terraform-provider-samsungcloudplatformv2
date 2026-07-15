provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_cloudmonitoring_account_events" "events" {
  event_state = var.EventState
  query_start_dt = var.QueryStartDt
  query_end_dt = var.QueryEndDt

}
