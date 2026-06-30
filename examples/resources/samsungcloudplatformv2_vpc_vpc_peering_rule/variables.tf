variable "vpc_peering_id" {
  type    = string
  default = "ENTER YOUR RESOURCE'S VPC_PEERING_ID"
}

variable "destination_cidr" {
  type    = string
  default = "77.33.0.0/24"
}

variable "destination_vpc_type" {
  type    = string
  default = "REQUESTER_VPC"
}


variable "tags" {
  type    = map(string)
  default = null
}


