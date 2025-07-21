variable "ip_address" {
  type = string
  default = "10.10.10.10"
}

variable "ip_id" {
  type = string
  default = "8a463aa4b1dc4f279c3f53b94dc45e74"
}

variable "ip_type" {
  type = string
  default = "PUBLIC"
}

variable "vpc_id" {
  type = string
  default = "8a463aa4b1dc4f279c3f53b94dc45e74"
}

variable "name" {
  type = string
  default = "terraformVpnGW"
}

variable "description" {
  type = string
  default = "description info"
}

variable "tags" {
  type    = map(string)
  default = {
    "vpn_tag_key" = "vpn_tag_value"
  }
}