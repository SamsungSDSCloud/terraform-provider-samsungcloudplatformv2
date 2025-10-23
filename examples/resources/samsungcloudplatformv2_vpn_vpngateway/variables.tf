variable "ip_address" {
  type    = string
  default = ""
}

variable "ip_id" {
  type    = string
  default = ""
}

variable "ip_type" {
  type    = string
  default = ""
}

variable "vpc_id" {
  type    = string
  default = ""
}

variable "name" {
  type    = string
  default = ""
}

variable "description" {
  type    = string
  default = ""
}

variable "tags" {
  type    = map(string)
  default = null
}

