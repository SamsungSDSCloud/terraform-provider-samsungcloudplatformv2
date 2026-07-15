variable "size" {
  type    = number
  default = 20
}

variable "page" {
  type    = number
  default = null
}

variable "sort" {
  type    = string
  default = null
}

variable "name" {
  type    = string
  default = null
}

variable "service_ip" {
  type    = string
  default = null
}

variable "subnet_id" {
  type    = string
  default = "ENTER YOUR RESOURCE'S SUBNET_ID"
}

variable "vpc_id" {
  type    = string
  default = "ENTER YOUR RESOURCE'S VPC_ID"
}



