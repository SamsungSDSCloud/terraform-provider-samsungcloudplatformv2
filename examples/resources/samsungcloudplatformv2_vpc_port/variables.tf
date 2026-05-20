variable "port_name" {
  type    = string
  default = ""
}

variable "subnet_id" {
  type    = string
  default = ""
}

variable "port_description" {
  type    = string
  default = ""
}

variable "port_fixed_ip_address" {
  type    = string
  default = ""
}

variable "security_groups" {
  type = list(object({
    id = string
  }))
  default = [{
    id = ""
  }]
}

