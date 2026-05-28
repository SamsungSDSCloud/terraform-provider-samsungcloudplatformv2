variable "firewall_connection_state" {
  type    = string
  default = null
}


variable "name" {
  type    = string
  default = null
}


variable "id" {
  type    = string
  default = "ENTER YOUR RESOURCE'S ID"
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

variable "state" {
  type    = string
  default = "ACTIVE"
}




