variable "port_name" {
  type    = string
  default = "testport"
}

variable "subnet_id" {
  type    = string
  default = "ENTER YOUR RESOURCE'S SUBNET_ID"
}

variable "port_description" {
  type    = string
  default = "test_description_sg"
}

variable "port_fixed_ip_address" {
  type    = string
  default = "192.168.0.13"
}

variable "security_groups" {
  type = list(object({
    id = string
  }))
  default = [{
    id = "ENTER YOUR RESOURCE'S ID"
  }]
}


