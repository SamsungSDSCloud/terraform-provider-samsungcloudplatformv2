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
  default = null
}

variable "vpn_gateway_id" {
  type    = string
  default = "ENTER YOUR RESOURCE'S VPN_GATEWAY_ID"
}

variable "vpn_gateway_name" {
  type    = string
  default = null
}




