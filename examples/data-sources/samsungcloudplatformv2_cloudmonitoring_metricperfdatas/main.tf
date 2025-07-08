provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_cloudmonitoring_metricperfdatas" "metric_perf_datas" {
  // header
  x_resource_type = var.XResourceType
  // request 필드가 계층구조더라도 테라폼에서는 필수값만 하나씩 입력 받기
  ignore_invalid = var.IgnoreInvalid
  metric_key = var.MetricKey
  product_resource_id = var.ProductResourceId
  query_end_dt = var.QueryEndDt
  query_start_dt = var.QueryStartDt
  object_list = var.ObjectList
  statistics_type_list = var.StatisticsTypeList
  object_type = var.ObjectType
  statistics_period  = var.StatisticsPeriod
}
// 계층구조
# data "samsungcloudplatformv2_cloudmonitoring_metricperfdatas" "metric_perf_datas" {
#   // header
#   x_resource_type = var.XResourceType
#   ignore_invalid = var.IgnoreInvalid
#   metric_data_conditions = var.MetricDataConditions
#   query_end_dt   = var.QueryEndDt
#   query_start_dt = var.QueryStartDt
# }
