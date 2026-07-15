variable "id" {
  type    = string
  default = "ENTER YOUR RESOURCE'S ID"
}

variable "name" {
  type    = string
  default = null
}

variable "vpc_id" {
  type    = string
  default = "ENTER YOUR RESOURCE'S VPC_ID"
}

variable "subnet_id" {
  type    = string
  default = "ENTER YOUR RESOURCE'S SUBNET_ID"
}

variable "vpc_name" {
  type    = string
  default = null
}

variable "resource_type" {
  type    = string
  default = null
}

variable "resource_key" {
  type    = string
  default = null
}

variable "endpoint_ip_address" {
  type    = string
  default = null
}

variable "state" {
  type    = string
  default = null
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


