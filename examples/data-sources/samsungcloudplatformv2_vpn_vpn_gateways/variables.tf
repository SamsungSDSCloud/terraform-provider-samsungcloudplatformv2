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

variable "ip_address" {
  type    = string
  default = null
}

variable "vpc_id" {
  type    = string
  default = "ENTER YOUR RESOURCE'S VPC_ID"
}

variable "vpc_name" {
  type    = string
  default = null
}



