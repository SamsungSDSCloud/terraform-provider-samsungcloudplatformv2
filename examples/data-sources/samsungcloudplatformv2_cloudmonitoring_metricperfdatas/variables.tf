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

