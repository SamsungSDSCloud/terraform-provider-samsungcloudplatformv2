variable "routing_rule_transit_gateway_id" {
  type    = string
  default = "ENTER YOUR RESOURCE'S ROUTING_RULE_TRANSIT_GATEWAY_ID"
}

variable "routing_rule_description" {
  type    = string
  default = "description VPC"
}

variable "routing_rule_destination_type" {
  type    = string
  default = "VPC"
}

variable "routing_rule_destination_cidr" {
  type    = string
  default = "1.1.0.0/17"
}

variable "routing_rule_tgw_connection_vpc_id" {
  type    = string
  default = "ENTER YOUR RESOURCE'S ROUTING_RULE_TGW_CONNECTION_VPC_ID"
}


