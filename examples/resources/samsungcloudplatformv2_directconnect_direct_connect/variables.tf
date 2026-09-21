variable "dcon_name" {
  type    = string
  default = "bactest2"
}

variable "dcon_vpc_id" {
  type    = string
  default = "ENTER YOUR RESOURCE'S DCON_VPC_ID"
}

variable "dcon_bandwidth" {
  type    = number
  default = 1
}

variable "dcon_description" {
  type    = string
  default = "description test"
}

variable "dcon_uplink_active_zone" {
  type    = string
  default = "kr-west1-a"
}

variable "dcon_uplink_standby_zone" {
  type    = string
  default = "kr-west1-y"
}

variable "dcon_firewall_enabled" {
  type    = bool
  default = false
}

variable "dcon_firewall_loggable" {
  type    = bool
  default = false
}

variable "tags" {
  type = map(string)
  default = {
    test_terraform_tag_key = "test_terraform_tag_value"
  }
}


