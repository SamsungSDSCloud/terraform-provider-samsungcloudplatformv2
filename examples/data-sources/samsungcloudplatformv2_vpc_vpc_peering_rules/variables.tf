variable "vpc_peering_id" {
  type    = string
  default = ""
}

variable "size" {
  type    = number
  default = 0
}

variable "page" {
  type    = number
  default = 0
}

variable "sort" {
  type    = string
  default = ""
}

variable "id" {
  type    = string
  default = ""
}

variable "source_vpc_id" {
  type    = string
  default = ""
}

variable "source_vpc_type" {
  type    = string
  default = ""
}

variable "destination_vpc_id" {
  type    = string
  default = ""
}

variable "destination_vpc_type" {
  type    = string
  default = ""
}

variable "destination_cidr" {
  type    = string
  default = ""
}

variable "state" {
  type    = string
  default = ""
}


