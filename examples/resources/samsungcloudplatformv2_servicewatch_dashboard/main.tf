provider "samsungcloudplatformv2" {
}

resource "samsungcloudplatformv2_servicewatch_dashboard" "dashboard" {
  name    = var.dashboard_name
  widgets = [
    {
      height = var.height
      width = var.width
      order = var.order
      type = "metric"
      properties = {
        period = var.period
        stacked = false
        statistic_type=var.statistic_type
        title = var.title
        view = "line"
        metrics = [{
            color=var.color
            display_name=var.display_name
            name=var.metric_name
            namespace_name=var.namespace_name
            period = var.period
            statistic_type = var.statistic_type
            dimensions = [
              {
                key=var.dimension_key
                value = var.dimension_value
              }
            ]
          }
        ]
      }
    }
  ]
}