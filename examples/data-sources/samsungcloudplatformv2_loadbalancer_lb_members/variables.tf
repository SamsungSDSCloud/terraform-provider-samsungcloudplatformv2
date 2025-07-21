variable "name" {
  type = string
  default = ""
}

variable "size" {
  type = number
  default = 20
}

variable "page" {
  type = number
  default = 0
}

variable "sort" {
  type = string
  default = ""
}

variable "lb_server_group_id" {
  type = string
  default = ""
}

variable "member_ip" {
  type = string
  default = "192.166.0.1"
}

variable "member_port" {
  type = number
  default = 80
}
