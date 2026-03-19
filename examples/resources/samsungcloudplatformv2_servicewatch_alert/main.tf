provider "samsungcloudplatformv2" {
}

resource "samsungcloudplatformv2_servicewatch_alert" "alert" {
  name                = var.alert_name
  description         = var.alert_description
  type                = var.alert_type
  level               = var.level
  activated_yn        = var.alert_activated_yn
  namespace_name      = var.namespace_name
  metric_name         = var.metric_name
  dimensions          = [{
      key               = var.dimension_key_1
      value             = var.dimension_value_1
    }, {
      key               = var.dimension_key_2
      value             = var.dimension_value_2
    }
  ]
  period              = var.period
  statistic           = var.statistic_type
  operator            = var.operator
  threshold           = var.threshold
  missing_data_option = var.missing_data_option
}