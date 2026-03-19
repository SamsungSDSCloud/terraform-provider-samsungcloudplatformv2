variable "event_ids" {
  type    = list(string)
  default = [""]
}

variable "event_rule_name" {
  type    = string
  default = ""
}

variable "recipient_ids" {
  type    = list(string)
  default = [""]
}

variable "service_id" {
  type    = string
  default = ""
}

variable "resource_type_id" {
  type    = string
  default = ""
}

variable "srn_list" {
  type    = list(string)
  default = [""]
}

variable "description" {
  type    = string
  default = ""
}

variable "tags" {
  type    = map(string)
  default = null
}

