variable "security_group_id" {
  type    = string
  default = ""
}

variable "ethertype" {
  type    = string
  default = ""
}

variable "protocol" {
  type    = string
  default = ""
}

variable "port_range_min" {
  type    = number
  default = 0
}

variable "port_range_max" {
  type    = number
  default = 0
}

variable "remote_ip_prefix" {
  type    = string
  default = ""
}

variable "remote_group_id" {
  type    = string
  default = ""
}


variable "description" {
  type    = string
  default = ""
}

variable "direction" {
  type    = string
  default = ""
}

