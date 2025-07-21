variable "id" {
  type = string
  default = ""
}

variable "name" {
  type = string
  default = ""
}

variable "ip" {
  type = string
  default = ""
}

variable "state" {
  type = string
  default = "ACTIVE"
}

variable "product_category" {
  type = string
  default = ""
}

variable "vpc_id" {
  type = string
  default = ""
}

variable "server_type_id" {
  type = string
  default = ""
}

variable "auto_scaling_group_id" {
  type = string
  default = ""
}

variable "server_filter_name" {
  type    = string
  default = ""
}

variable "server_filter_values" {
  type    = list(string)
  default = [""]
}

variable "server_filter_use_regex" {
  type    = bool
  default = true
}