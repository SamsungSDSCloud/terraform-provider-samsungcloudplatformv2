
variable "cidr" {
  type    = string
  default = ""
}

variable "id" {
  type    = string
  default = ""
}

variable "name" {
  type    = string
  default = ""
}

variable "page" {
  type    = number
  default = 0
}

variable "size" {
  type    = number
  default = 0
}

variable "sort" {
  type    = string
  default = ""
}

variable "state" {
  type    = string
  default = ""
}

variable "type" {
  type    = list(string)
  default = [""]
}

variable "vpc_id" {
  type    = string
  default = ""
}

variable "vpc_name" {
  type    = string
  default = ""
}


