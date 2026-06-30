variable "igw_vpc_id" {
  type    = string
  default = "ENTER YOUR RESOURCE'S IGW_VPC_ID"
}

variable "igw_type" {
  type    = string
  default = "GGW"
}

variable "igw_description" {
  type    = string
  default = "igw update test 4"
}

variable "loggable" {
  type    = bool
  default = null
}

variable "firewall_enabled" {
  type    = bool
  default = null
}

variable "firewall_loggable" {
  type    = bool
  default = null
}

variable "tags" {
  type = map(string)
  default = {
    test_tag_key = "test_tag_value"
  }
}


