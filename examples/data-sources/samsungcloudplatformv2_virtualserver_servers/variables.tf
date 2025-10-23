variable "name" {
  type    = string
  default = ""
}

variable "ip" {
  type    = string
  default = ""
}

variable "state" {
  type    = string
  default = ""
}

variable "product_category" {
  type    = string
  default = ""
}

variable "vpc_id" {
  type    = string
  default = ""
}

variable "server_type_id" {
  type    = string
  default = ""
}

variable "auto_scaling_group_id" {
  type    = string
  default = ""
}

variable "servers_filter_name" {
  type    = string
  default = ""
}

variable "servers_filter_values" {
  type    = list(string)
  default = [""]
}

variable "servers_filter_use_regex" {
  type    = bool
  default = false
}

