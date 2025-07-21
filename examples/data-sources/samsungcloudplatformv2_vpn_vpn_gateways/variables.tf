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
  default = ""
}

variable "ip_address" {
  type    = string
  default = ""
}

variable "vpc_id" {
  type    = string
  default = ""
}

variable "vpc_name" {
  type    = string
  default = ""
}
