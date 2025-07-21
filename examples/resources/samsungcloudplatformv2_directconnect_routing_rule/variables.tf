variable "routing_rule_direct_connect_id" {
  type = string
  default = "7df8abb4912e4709b1cb237daccca7a8"
}

variable "routing_rule_destination_type" {
  type = string
  default = "ON_PREMISE"
}

variable "routing_rule_destination_cidr" {
  type = string
  default = "10.10.10.0/24"
}