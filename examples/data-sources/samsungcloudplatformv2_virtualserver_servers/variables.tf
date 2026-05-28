variable "name" {
  type    = string
  default = null
}

variable "ip" {
  type    = string
  default = null
}

variable "state" {
  type    = string
  default = "ACTIVE"
}

variable "product_category" {
  type    = string
  default = null
}

variable "vpc_id" {
  type    = string
  default = "ENTER YOUR RESOURCE'S VPC_ID"
}

variable "server_type_id" {
  type    = string
  default = "ENTER YOUR RESOURCE'S SERVER_TYPE_ID"
}

variable "auto_scaling_group_id" {
  type    = string
  default = "ENTER YOUR RESOURCE'S AUTO_SCALING_GROUP_ID"
}

variable "servers_filter_name" {
  type    = string
  default = "name"
}

variable "servers_filter_values" {
  type    = list(string)
  default = ["test"]
}

variable "servers_filter_use_regex" {
  type    = bool
  default = true
}


