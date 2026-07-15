variable "cluster_cloud_logging_enabled" {
  type    = bool
  default = true
}

variable "cluster_kubernetes_version" {
  type    = string
  default = "v1.31.8"
}

variable "cluster_name" {
  type    = string
  default = "terraform-test"
}

variable "cluster_security_group_id_list" {
  type    = list(string)
  default = ["852ecca3-7252-44a9-9582-967745ae338c", "09029278-adca-4acf-b959-4eabc3d4baeb"]
}

variable "cluster_subnet_id" {
  type    = string
  default = "ENTER YOUR RESOURCE'S CLUSTER_SUBNET_ID"
}

variable "cluster_vpc_id" {
  type    = string
  default = "ENTER YOUR RESOURCE'S CLUSTER_VPC_ID"
}

variable "cluster_volume_id" {
  type    = string
  default = "ENTER YOUR RESOURCE'S CLUSTER_VOLUME_ID"
}

variable "private_endpoint_access_control_resources" {
  type = list(object({
    id   = string
    name = string
    type = string
  }))
  default = []
}

variable "public_endpoint_access_control_ip" {
  type    = string
  default = "123.123.123.123"
}

variable "service_watch_logging_enabled" {
  type    = bool
  default = false
}


