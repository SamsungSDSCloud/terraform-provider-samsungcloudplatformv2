variable "EventLevel" {
  type    = string
  default = ""
}

variable "MetricKey" {
  type    = string
  default = ""
}

variable "DisableYn" {
  type    = string
  default = ""
}

variable "EventMessagePrefix" {
  type    = string
  default = ""
}

variable "FtCount" {
  type = number

  default = 0
}

variable "IsLogMetric" {
  type    = string
  default = ""
}

variable "ObjectName" {
  type    = string
  default = ""
}

variable "ObjectDisplayName" {
  type    = string
  default = ""
}

variable "EventOccurTimeZone" {
  type    = string
  default = ""
}

variable "ObjectType" {
  type    = string
  default = ""
}

variable "ObjectTypeName" {
  type    = string
  default = ""
}



variable "EventState" {
  type    = string
  default = ""
}

variable "QueryStartDt" {
  type    = string
  default = ""
}

variable "QueryEndDt" {
  type    = string
  default = ""
}

variable "XResourceType" {
  type    = string
  default = ""
}

variable "EventId" {
  type    = string
  default = ""
}

variable "ProductResourceId" {
  type    = string
  default = ""
}

variable "EventPolicyId" {
  type = number


  default = 0
}

variable "EventThreshold" {
  type = object({
    event_threshold = object({
      metric_function = string
      single_threshold = object({
        comparison_operator = string
        value               = number
      })
      threshold_type = string
    })
  })

  default = {
    event_threshold = {
      metric_function = ""
      single_threshold = {
        comparison_operator = ""
        value               = 0
      }
      threshold_type = ""
    }
  }
}

variable "EventPolicy" {
  type = object({
    asg_yn = string
    #    product_resource_id     = string
    event_level          = string
    disable_yn           = string
    start_dt             = string
    event_message_prefix = string
    ft_count             = number
    is_log_metric        = string
    metric_key           = string
    #    object_name             = string
    object_type         = string
    object_type_name    = string
    object_display_name = string
    disable_object      = string
    attr_list_str       = string
    #        pod_object_name         = ""
    pod_object_display_name = string
    event_rule              = string
    display_event_rule      = string
    #        metric_name             = ""
    event_occur_time_zone = string
    #        object_type             = ""
    #        object_type_name        = ""
    metric_description    = string
    metric_description_en = string
    user_names            = string
    user_name_str         = string
    event_threshold = object({
      metric_function = string
      single_threshold = object({
        comparison_operator = string
        value               = number
      })
      threshold_type = string
    })
  })

  default = {
    asg_yn                = ""
    attr_list_str         = ""
    disable_object        = ""
    disable_yn            = ""
    display_event_rule    = ""
    event_level           = ""
    event_message_prefix  = ""
    event_occur_time_zone = ""
    event_rule            = ""
    event_threshold = {
      metric_function = ""
      single_threshold = {
        comparison_operator = ""
        value               = 0
      }
      threshold_type = ""
    }
    ft_count                = 0
    is_log_metric           = ""
    metric_description      = ""
    metric_description_en   = ""
    metric_key              = ""
    object_display_name     = ""
    object_type             = ""
    object_type_name        = ""
    pod_object_display_name = ""
    start_dt                = ""
    user_name_str           = ""
    user_names              = ""
  }
}

