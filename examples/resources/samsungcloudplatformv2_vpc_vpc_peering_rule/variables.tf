variable "vpc_peering_id" {
  type    = string
  default = ""
}

variable "destination_cidr" {
  type    = string
  default = ""
}

variable "destination_vpc_type" {
  type    = string
  default = ""
}


variable "tags" {
  type    = map(string)
  default = null
}

