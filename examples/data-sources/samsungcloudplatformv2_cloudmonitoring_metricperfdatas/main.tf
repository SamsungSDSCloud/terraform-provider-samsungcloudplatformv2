provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_cloudmonitoring_metricperfdatas" "metric_perf_datas" {
  // header
  x_resource_type = var.XResourceType
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