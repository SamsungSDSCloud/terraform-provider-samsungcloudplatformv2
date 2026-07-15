variable "vpcendpoint_name" {
  type    = string
  default = "sample"
}

variable "vpcendpoint_vpc_id" {
  type    = string
  default = "ENTER YOUR RESOURCE'S VPCENDPOINT_VPC_ID"
}

variable "vpcendpoint_subnet_id" {
  type    = string
  default = "ENTER YOUR RESOURCE'S VPCENDPOINT_SUBNET_ID"
}

variable "vpcendpoint_resource_type" {
  type    = string
  default = "sample-resource-type"
}

variable "vpcendpoint_resource_resource_key" {
  type    = string
  default = "sample-resource-key"
}

variable "vpcendpoint_resource_resource_info" {
  type    = string
  default = "sample-resource-info"
}

variable "vpcendpoint_endpoint_ip_address" {
  type    = string
  default = "sample-endpoint-ip-address"
}

variable "vpcendpoint_description" {
  type    = string
  default = "sample-description"
}

variable "vpcendpoint_tags" {
  type = map(string)
  default = {
    tf_key1 = "tf_val1"
  }
}



