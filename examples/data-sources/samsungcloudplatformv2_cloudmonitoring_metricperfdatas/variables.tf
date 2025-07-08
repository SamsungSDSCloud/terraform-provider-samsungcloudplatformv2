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

