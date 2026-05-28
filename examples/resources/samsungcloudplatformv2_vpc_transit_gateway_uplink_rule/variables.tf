variable "description" {
  type    = string
  default = "TGW cidr one of vpc"
}

variable "destination_cidr" {
  type    = string
  default = "1.1.0.0/17"
}

variable "destination_type" {
  type    = string
  default = "TGW"
}

variable "transit_gateway_id" {
  type    = string
  default = "ENTER YOUR RESOURCE'S TRANSIT_GATEWAY_ID"
}


