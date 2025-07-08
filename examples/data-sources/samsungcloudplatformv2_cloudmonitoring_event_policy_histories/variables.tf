
## qa2

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

variable "EventPolicyId" {
  type    = number
  default = 0
}


## dev2
#
#variable "QueryStartDt" {
#  type = string
#  default = "2025-03-01T23:00:00.000Z"
#}
#
#variable "QueryEndDt" {
#  type = string
#  default = "2025-04-02T08:00:00.000Z"
#}
#
#variable "XResourceType" {
#  type = string
#  default = "VM"
#}
#
#variable "EventPolicyId" {
#  type = number
#  # default = 30878       ##  dev2
#  default = 12195         ##  qa2
#}


