variable "QueryStartDt" {
  type = string
  default = "2025-07-14T00:00:00.000Z"
}

variable "QueryEndDt" {
  type = string
  default = "2025-07-17T08:00:00.000Z"
}

variable "XResourceType" {
  type = string
  default = "VM"
}

variable "EventPolicyId" {
  type = number
  default = 30878
}