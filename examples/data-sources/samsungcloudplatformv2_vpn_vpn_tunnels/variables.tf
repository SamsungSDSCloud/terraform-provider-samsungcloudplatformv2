variable "size" {
  type    = number
  default = 20
}

variable "page" {
  type    = number
  default = 0
}

variable "sort" {
  type    = string
  default = "created_at:desc"
}

variable "name" {
  type    = string
  default = ""
}

variable "vpn_gateway_id" {
  type    = string
  default = ""
}

variable "vpn_gateway_name" {
  type    = string
  default = ""
}

variable "peer_gateway_ip" {
  type    = string
  default = ""
}

variable "remote_subnet" {
  type    = string
  default = ""
}
