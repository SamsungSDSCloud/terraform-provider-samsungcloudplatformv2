variable "subnet_id" {
  type    = string
  default = "ENTER YOUR RESOURCE'S SUBNET_ID"
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
  default = "created_at:desc"
}

variable "virtual_ip_address" {
  type    = string
  default = null
}

variable "public_ip_address" {
  type    = string
  default = null
}



