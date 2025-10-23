variable "region" {
  type    = string
  default = ""
}

variable "name" {
  type    = string
  default = ""
}

variable "policy_category" {
  type    = string
  default = ""
}

variable "policy_type" {
  type    = string
  default = ""
}

variable "server_uuid" {
  type    = string
  default = ""
}

variable "server_category" {
  type    = string
  default = ""
}

variable "encrypt_enabled" {
  type    = bool
  default = false
}

variable "retention_period" {
  type    = string
  default = ""
}

variable "schedules" {
  type = list(object({
    type       = string
    frequency  = string
    start_time = string
    start_day  = string
    start_week = string
  }))
  default = [{
    frequency  = ""
    start_day  = ""
    start_time = ""
    start_week = ""
    type       = ""
  }]
}

variable "tags" {
  type    = map(string)
  default = null
}


