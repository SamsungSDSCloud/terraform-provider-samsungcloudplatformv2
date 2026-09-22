variable "natgateway_subnet_id" {
  type    = string
  default = "ENTER YOUR RESOURCE'S NATGATEWAY_SUBNET_ID"
}

variable "natgateway_publicip_ids" {
  type    = list(string)
  default = ["ENTER YOUR RESOURCE'S NATGATEWAY_PUBLICIP_IDS"]
}

variable "natgateway_multi_zone_enabled" {
  type    = bool
  default = false
}

variable "natgateway_description" {
  type    = string
  default = "description-by-terrafrom"
}



