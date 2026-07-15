variable "XResourceType" {
  type    = string
  default = "PostgreSQL"
}
variable "IgnoreInvalid" {
  type    = string
  default = "Y"
}
variable "QueryEndDt" {
  type    = string
  default = "2025-06-19T07:55:39.569Z"
}
variable "QueryStartDt" {
  type    = string
  default = "2025-06-19T07:25:39.569Z"
}
variable "MetricKey" {
  type    = string
  default = "postgresql.status.pid"
}
variable "ProductResourceId" {
  type    = string
  default = "ENTER YOUR RESOURCE'S PRODUCTRESOURCEID"
}
variable "ObjectList" {
  type    = list(string)
  default = null
}
variable "ObjectType" {
  type    = string
  default = null
}
variable "StatisticsTypeList" {
  type    = list(string)
  default = null
}
variable "StatisticsPeriod" {
  type    = number
  default = null
}


