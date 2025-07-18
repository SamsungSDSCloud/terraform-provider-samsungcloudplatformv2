---
page_title: "samsungcloudplatformv2_cloudmonitoring_metricperfdatas Data Source - samsungcloudplatformv2"
subcategory: samsungcloudplatformv2_cloudmonitoring_metricperfdatas
description: |-
  The Schema of MetricPerfDataDataSources.
---

# samsungcloudplatformv2_cloudmonitoring_metricperfdatas (Data Source)

The Schema of MetricPerfDataDataSources.

## Example Usage

```terraform
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


output "metric_perf_datas" {
  value = data.samsungcloudplatformv2_cloudmonitoring_metricperfdatas.metric_perf_datas
}

variable "XResourceType" {
  type    = string
  default = ""
}
variable "IgnoreInvalid" {
  type    = string
  default = ""
}
variable "QueryEndDt" {
  type    = string
  default = ""
}
variable "QueryStartDt" {
  type    = string
  default = ""
}
# variable "MetricDataConditions" {
#   type    = list(object({
#     metricKey             = string
#     productResourceInfos  = list(object({
#       productResourceId = string
#     }))
#     statisticsPeriod      = number
#     statisticsTypeList    = list(string)
#   }))
#   default = [
#     {
#       metricKey             = "postgresql.status.pid"
#       productResourceInfos  = [
#         {
#           productResourceId = "36d845e6-9f26-4f35-9e2a-0a9456858742"
#         }
#       ]
#       statisticsPeriod      = null
#       statisticsTypeList    = null
#     }
#   ]
# }
variable "MetricKey" {
  type    = string
  default = ""
}
variable "ProductResourceId" {
  type    = string
  default = ""
}
variable "ObjectList" {
  type    = list(string)
  default = [""]
}
variable "ObjectType" {
  type    = string
  default = ""
}
variable "StatisticsTypeList" {
  type    = list(string)
  default = [""]
}
variable "StatisticsPeriod" {
  type    = number
  default = 0
}
// body
# variable "request" {
#   type = string
#   default = "{\"ignoreInvalid\": \"Y\",\"metricDataConditions\": [  {\"metricKey\": \"postgresql.status.pid\",\"objectType\": \"\",\"productResourceInfos\": [  {\"objectList\": [],\"productResourceId\": \"36d845e6-9f26-4f35-9e2a-0a9456858742\"  }],\"statisticsPeriod\": null,\"statisticsTypeList\": null  }],\"queryEndDt\": \"2025-06-19T07:55:39.569Z\",\"queryStartDt\": \"2025-06-19T07:25:39.569Z\"  }"
# }
# variable "Request" {
#   type = object({
#     ignoreInvalid = string
#     metricDataConditions = list(object({
#       metricKey           = string
#       objectType          = string
#       productResourceInfos = list(object({
#         objectList        = list(string)
#         productResourceId = string
#       }))
#       statisticsPeriod = number
#       statisticsTypeList = list(string)
#     }))
#     queryEndDt   = string
#     queryStartDt = string
#   })
#   default = {
#     ignoreInvalid = "Y"
#     metricDataConditions = [
#       {
#         metricKey           = "postgresql.status.pid"
#         objectType          = ""
#         productResourceInfos = [
#           {
#             objectList        = []
#             productResourceId = "36d845e6-9f26-4f35-9e2a-0a9456858742"
#           }
#         ]
#         statisticsPeriod = null
#         statisticsTypeList = null
#       }
#     ]
#     queryEndDt   = "2025-06-19T07:55:39.569Z"
#     queryStartDt = "2025-06-19T07:25:39.569Z"
#   }
# }
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `metric_key` (String) MetricKey
- `product_resource_id` (String) ProductResourceId
- `query_end_dt` (String) QueryEndDt
- `query_start_dt` (String) QueryStartDt
- `x_resource_type` (String) X-ResourceType

### Optional

- `filter` (Block List) Filter (see [below for nested schema](#nestedblock--filter))
- `ignore_invalid` (String) IgnoreInvalid
- `object_list` (List of String) ObjectList
- `object_type` (String) ObjectType
- `statistics_period` (Number) StatisticsPeriod
- `statistics_type_list` (List of String) StatisticsTypeList

### Read-Only

- `metric_perf_datas` (Attributes List) MetricPerfDatas (see [below for nested schema](#nestedatt--metric_perf_datas))

<a id="nestedblock--filter"></a>
### Nested Schema for `filter`

Required:

- `name` (String) Filtering target name
- `use_regex` (Boolean) Enable regex match for values
- `values` (List of String) Filtering values. Each matching value is appended. (OR rule)


<a id="nestedatt--metric_perf_datas"></a>
### Nested Schema for `metric_perf_datas`

Required:

- `statistics_period` (Number) StatisticsPeriod
- `statistics_type` (String) StatisticsType

Optional:

- `metric_key` (String) MetricKey
- `metric_name` (String) MetricName
- `metric_type` (String) MetricType
- `metric_unit` (String) MetricUnit
- `object_display_name` (String) ObjectDisplayName
- `object_name` (String) ObjectName
- `object_type` (String) ObjectType
- `product_name` (String) ProductName
- `product_resource_id` (String) ProductResourceId

Read-Only:

- `perf_data` (Attributes List) PerfData (see [below for nested schema](#nestedatt--metric_perf_datas--perf_data))

<a id="nestedatt--metric_perf_datas--perf_data"></a>
### Nested Schema for `metric_perf_datas.perf_data`

Optional:

- `ts` (Number) ts
- `value` (String) value
