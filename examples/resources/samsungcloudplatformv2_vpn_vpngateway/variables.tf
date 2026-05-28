variable "ip_address" {
  type    = string
  default = "42.14.6.17"
}

variable "ip_id" {
  type    = string
  default = "ENTER YOUR RESOURCE'S IP_ID"
}

variable "ip_type" {
  type    = string
  default = "PUBLIC"
}

variable "vpc_id" {
  type    = string
  default = "ENTER YOUR RESOURCE'S VPC_ID"
}

variable "name" {
  type    = string
  default = "terraformVpnGW"
}

variable "description" {
  type    = string
  default = "test_123"
}

variable "tags" {
  type = map(string)
  default = {
    vpn_tag_key = "vpn_tag_value"
  }
}


