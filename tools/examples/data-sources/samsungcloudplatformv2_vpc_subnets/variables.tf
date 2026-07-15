
variable "cidr" {
  type    = string
  default = null
}

variable "id" {
  type    = string
  default = "ENTER YOUR RESOURCE'S ID"
}

variable "name" {
  type    = string
  default = null
}

variable "page" {
  type    = number
  default = 0
}

variable "size" {
  type    = number
  default = 20
}

variable "sort" {
  type    = string
  default = "created_at:desc"
}

variable "state" {
  type    = string
  default = null
}

variable "type" {
  type    = list(string)
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



