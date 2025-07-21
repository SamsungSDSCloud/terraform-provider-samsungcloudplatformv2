variable "XResourceType" {
  type = string
  default = "PostgreSQL"
}
variable "IgnoreInvalid" {
  type = string
  default = "Y"
}
variable "QueryEndDt" {
  type = string
  default = "2025-07-19T07:55:39.569Z"
}
variable "QueryStartDt" {
  type = string
  default = "2025-07-19T07:25:39.569Z"
}

variable "MetricKey" {
  type = string
  default = "postgresql.status.pid"
}
variable "ProductResourceId" {
  type = string
  default = "8a463aa4-b1dc-4f27-9c3f-53b94dc45e74"
}
variable "ObjectList" {
  type = list(string)
  default = [""]
}
variable "ObjectType" {
  type = string
  default = ""
}
variable "StatisticsTypeList" {
  type = list(string)
  default = [""]
}
variable "StatisticsPeriod" {
  type = number
  default = null
}
