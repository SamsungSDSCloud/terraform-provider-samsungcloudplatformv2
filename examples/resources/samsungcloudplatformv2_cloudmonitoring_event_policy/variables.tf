variable "EventLevel" {
  type    = string
  default = "FATAL"
}

variable "MetricKey" {
  type    = string
  default = "libvirt.domain.cpu.scpm.usage"
}

variable "DisableYn" {
  type    = string
  default = "N"
}

variable "EventMessagePrefix" {
  type    = string
  default = ""
}

variable "FtCount" {
  type    = number
  default = 1

}

variable "IsLogMetric" {
  type    = string
  default = "N"
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
  default = "GMT+0"
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
  default = "ALL"
}

variable "QueryStartDt" {
  type    = string
  default = "2025-01-01T18:00:00.000Z"
}

variable "QueryEndDt" {
  type    = string
  default = "2025-03-12T23:59:59.000Z"
}

variable "XResourceType" {
  type    = string
  default = "MS-SQL"
}

variable "EventId" {
  type    = string
  default = "ENTER YOUR RESOURCE'S EVENTID"
}

variable "ProductResourceId" {
  type    = string
  default = "ENTER YOUR RESOURCE'S PRODUCTRESOURCEID"
}

variable "EventPolicyId" {
  type    = number
  default = 67762


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
      metric_function = "delta"
      single_threshold = {
        comparison_operator = "GE"
        value               = 75
      }
      threshold_type = "TEST_SINGLE_VALUE"
    }
  }
}

variable "EventPolicy" {
  type = object({
    asg_yn                  = string
    event_level             = string
    disable_yn              = string
    start_dt                = string
    event_message_prefix    = string
    ft_count                = number
    is_log_metric           = string
    metric_key              = string
    object_type             = string
    object_type_name        = string
    object_display_name     = string
    disable_object          = string
    attr_list_str           = string
    pod_object_display_name = string
    event_rule              = string
    display_event_rule      = string
    event_occur_time_zone   = string
    metric_description      = string
    metric_description_en   = string
    user_names              = string
    user_name_str           = string
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
    asg_yn                = "N"
    attr_list_str         = ""
    disable_object        = ""
    disable_yn            = "N"
    display_event_rule    = "del(CPU Usage [Basic]) >= 75.0"
    event_level           = "FATAL"
    event_message_prefix  = ""
    event_occur_time_zone = "GMT+0"
    event_rule            = ""
    event_threshold = {
      metric_function = "delta"
      single_threshold = {
        comparison_operator = "GE"
        value               = 75
      }
      threshold_type = "SINGLE_VALUE"
    }
    ft_count                = 1
    is_log_metric           = "N"
    metric_description      = "Idle 및 IOWait 상태 이외에 사용된 CPU 시간의 백분율"
    metric_description_en   = "Percentage of CPU time used outside of Idle and IOWait states"
    metric_key              = "libvirt.domain.cpu.scpm.usage"
    object_display_name     = ""
    object_type             = ""
    object_type_name        = ""
    pod_object_display_name = ""
    start_dt                = "0001-01-01 00:00:00 +0000 UTC"
    user_name_str           = ""
    user_names              = ""
  }
}


