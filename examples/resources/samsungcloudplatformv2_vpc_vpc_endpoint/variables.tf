variable "vpcendpoint_name" {
  type = string
  default = "vpcEndpointName"
}

variable "vpcendpoint_vpc_id" {
  type = string
  default = "7df8abb4912e4709b1cb237daccca7a8"
}

variable "vpcendpoint_subnet_id" {
  type = string
  default = "7df8abb4912e4709b1cb237daccca7a8"
}

variable "vpcendpoint_resource_type" {
  type = string
  default = "FS"
}

variable "vpcendpoint_resource_resource_key" {
  type = string
  default = "1.1.1.1"
}

variable "vpcendpoint_resource_resource_info" {
  type = string
  default = "192.168.0.1(SSD)"
}

variable "vpcendpoint_endpoint_ip_address" {
  type = string
  default = "10.10.10.10"
}

variable "vpcendpoint_description" {
  type = string
  default = "description info"
}

variable "vpcendpoint_tags" {
  type    = map(string)
  default = {
    "tf_key1" = "tf_val1"
  }
}
