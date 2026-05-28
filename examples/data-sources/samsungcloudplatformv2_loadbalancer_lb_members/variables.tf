variable "name" {
  type    = string
  default = null
}

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

variable "lb_server_group_id" {
  type    = string
  default = "ENTER YOUR RESOURCE'S LB_SERVER_GROUP_ID"
}

variable "member_ip" {
  type    = string
  default = null
}

variable "member_port" {
  type    = number
  default = null
}



