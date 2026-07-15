variable "transit_gateway_id" {
  type    = string
  default = "ENTER YOUR RESOURCE'S TRANSIT_GATEWAY_ID"
}

variable "size" {
  type    = number
  default = 10
}

variable "page" {
  type    = number
  default = 0
}

variable "sort" {
  type    = string
  default = "created_at:desc"
}

variable "id" {
  type    = string
  default = "ENTER YOUR RESOURCE'S ID"
}

variable "tgw_connection_vpc_id" {
  type    = string
  default = "ENTER YOUR RESOURCE'S TGW_CONNECTION_VPC_ID"
}

variable "tgw_connection_vpc_name" {
  type    = string
  default = "test0819"
}

variable "source_type" {
  type    = string
  default = "TGW"
}

variable "destination_type" {
  type    = string
  default = "VPC"
}

variable "destination_cidr" {
  type    = string
  default = "1.1.1.0/24"
}

variable "state" {
  type    = string
  default = "ACTIVE"
}

variable "rule_type" {
  type    = string
  default = null
}




