variable "security_group_id" {
  type = string
  default = "8a463aa4-b1dc-4f27-9c3f-53b94dc45e74"
}

variable "ethertype" {
  type = string
  default = "IPv4"
}

variable "protocol" {
  type = string
  default = "TCP"
}

variable "port_range_min" {
  type = number
  default = 22
}

variable "port_range_max" {
  type = number
  default = 23
}

variable "remote_ip_prefix" {
  type = string
  default = "1.1.1.1/30"
}

variable "remote_group_id" {
  type = string
  default = ""
}


variable "description" {
  type = string
  default = "description info"
}

variable "direction" {
  type = string
  default = "egress"
}