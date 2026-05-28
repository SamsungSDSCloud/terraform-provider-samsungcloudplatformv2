variable "private_nat_cidr" {
  type    = string
  default = "192.167.0.0/24"
}

variable "private_nat_name" {
  type    = string
  default = "sdsv-test-01"
}

variable "private_nat_service_resource_id" {
  type    = string
  default = "ENTER YOUR RESOURCE'S PRIVATE_NAT_SERVICE_RESOURCE_ID"
}

variable "private_nat_service_type" {
  type    = string
  default = "TRANSIT_GATEWAY"
}

variable "private_nat_description" {
  type    = string
  default = "description upda"
}

variable "private_nat_tags" {
  type = map(string)
  default = {
    tags1 = "tag_me"
  }
}



