variable "region" {
  type    = string
  default = ""
}

variable "server_name" {
  type    = string
  default = ""
}

variable "state" {
  type    = string
  default = ""
}

variable "vpc_id" {
  type    = string
  default = ""
}

variable "ip" {
  type    = string
  default = ""
}

variable "baremetals_filter_name" {
  type    = string
  default = ""
}

variable "baremetals_filter_values" {
  type    = list(string)
  default = [""]
}

variable "baremetals_filter_use_regex" {
  type    = bool
  default = false
}

