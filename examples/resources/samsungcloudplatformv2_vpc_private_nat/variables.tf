variable "private_nat_name" {
  type    = string
  default = ""
}

variable "private_nat_direct_connect_id" {
  type    = string
  default = ""
}

variable "private_nat_cidr" {
  type    = string
  default = ""
}

variable "private_nat_description" {
  type    = string
  default = ""
}

variable "private_nat_tags" {
  type    = map(string)
  default = null
}


