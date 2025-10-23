variable "cluster_cloud_logging_enabled" {
  type    = bool
  default = false
}

variable "cluster_kubernetes_version" {
  type    = string
  default = ""
}

variable "cluster_name" {
  type    = string
  default = ""
}

variable "cluster_security_group_id_list" {
  type    = list(string)
  default = [""]
}

variable "cluster_subnet_id" {
  type    = string
  default = ""
}

variable "cluster_vpc_id" {
  type    = string
  default = ""
}

variable "cluster_volume_id" {
  type    = string
  default = ""
}

variable "private_endpoint_access_control_resources" {
  type = list(object({
    id   = string
    name = string
    type = string
  }))
  default = [{
    id   = ""
    name = ""
    type = ""
  }]
}

variable "public_endpoint_access_control_ip" {
  type    = string
  default = ""
}

variable "service_watch_logging_enabled" {
  type    = bool
  default = false
}

