variable "cn" {
  type    = string
  default = ""
}

variable "not_after_dt" {
  type    = string
  default = ""
}

variable "not_before_dt" {
  type    = string
  default = ""
}

variable "name" {
  type    = string
  default = ""
}

variable "organization" {
  type    = string
  default = ""
}

variable "region" {
  type    = string
  default = ""
}

variable "timezone" {
  type    = string
  default = ""
}

variable "recipients" {
  type    = list(map(string))
  default = [null]
}

variable "tags" {
  type    = map(string)
  default = null
}

