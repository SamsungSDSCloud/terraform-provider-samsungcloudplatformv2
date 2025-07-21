variable "security_group_id" {
  type = string
  default =  "8a463aa4-b1dc-4f27-9c3f-53b94dc45e74"
}

variable "size" {
  type    = number
  default = 10
}

variable "page" {
  type    = number
  default = 0
}

variable "id" {
  type    = string
  default = ""
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

variable "service" {
  type    = string
  default = ""
}