variable "transit_gateway_id" {
  type    = string
  default = "ENTER YOUR RESOURCE'S TRANSIT_GATEWAY_ID"
}

variable "product_type" {
  type    = string
  default = "TGW_DGW"
}

variable "uplink_active_zone" {
  type    = string
  default = "kr-west1-a"
}

variable "uplink_standby_zone" {
  type    = string
  default = "kr-west1-y"
}



