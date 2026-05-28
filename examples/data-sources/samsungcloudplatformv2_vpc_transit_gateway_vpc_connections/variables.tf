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


variable "transit_gateway_id" {
  type    = string
  default = "ENTER YOUR RESOURCE'S TRANSIT_GATEWAY_ID"
}

variable "id" {
  type    = string
  default = "ENTER YOUR RESOURCE'S ID"
}

variable "vpc_id" {
  type    = string
  default = "ENTER YOUR RESOURCE'S VPC_ID"
}

variable "vpc_name" {
  type    = string
  default = null
}

variable "state" {
  type    = string
  default = "ACTIVE"
}




