variable "direct_connect_id" {
  type    = string
  default = "ENTER YOUR RESOURCE'S DIRECT_CONNECT_ID"
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

variable "id" {
  type    = string
  default = "ENTER YOUR RESOURCE'S ID"
}

variable "destination_type" {
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



