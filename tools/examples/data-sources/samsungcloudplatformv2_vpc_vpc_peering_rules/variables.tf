variable "vpc_peering_id" {
  type    = string
  default = "ENTER YOUR RESOURCE'S VPC_PEERING_ID"
}

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
  default = null
}

variable "id" {
  type    = string
  default = "ENTER YOUR RESOURCE'S ID"
}

variable "source_vpc_id" {
  type    = string
  default = "ENTER YOUR RESOURCE'S SOURCE_VPC_ID"
}

variable "source_vpc_type" {
  type    = string
  default = null
}

variable "destination_vpc_id" {
  type    = string
  default = "ENTER YOUR RESOURCE'S DESTINATION_VPC_ID"
}

variable "destination_vpc_type" {
  type    = string
  default = null
}

variable "destination_cidr" {
  type    = string
  default = null
}

variable "state" {
  type    = string
  default = null
}



