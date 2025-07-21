variable "port_name" {
  type = string
  default = "portName"
}

variable "subnet_id" {
  type = string
  default = "023c57b14f11483689338d085e061492"
}

variable "port_description" {
  type = string
  default = "description info"
}

variable "port_fixed_ip_address" {
    type = string
    default = "10.10.10.10"
}

variable "port_security_groups" {
  type = list(string)
  default = [
    "3eef50bc-d638-41fa-99f3-5f9a877dd864",
    "b81d2ec8-b896-4853-bc7d-b06a5f28e228"
  ]
}