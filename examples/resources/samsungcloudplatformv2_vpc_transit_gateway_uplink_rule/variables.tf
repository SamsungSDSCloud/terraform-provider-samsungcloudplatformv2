variable "description" {
  type    = string
  default = "TGW rule nam dep trai"
}

variable "destination_cidr" {
  type    = string
  default = "192.168.100.0/24"
}

variable "destination_type" {
  type    = string
  default = "TGW"
}

variable "transit_gateway_id" {
  type    = string
  default = "ENTER YOUR RESOURCE'S TRANSIT_GATEWAY_ID"
}


