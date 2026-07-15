variable "security_group_id" {
  type    = string
  default = "ENTER YOUR RESOURCE'S SECURITY_GROUP_ID"
}

variable "ethertype" {
  type    = string
  default = "IPv4"
}

variable "protocol" {
  type    = string
  default = "TCP"
}

variable "port_range_min" {
  type    = number
  default = 22
}

variable "port_range_max" {
  type    = number
  default = 23
}

variable "remote_ip_prefix" {
  type    = string
  default = "1.1.1.1/30"
}

variable "remote_group_id" {
  type    = string
  default = "ENTER YOUR RESOURCE'S REMOTE_GROUP_ID"
}


variable "description" {
  type    = string
  default = "test description"
}

variable "direction" {
  type    = string
  default = "egress"
}


