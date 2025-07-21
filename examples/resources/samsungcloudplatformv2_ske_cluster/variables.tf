variable "cluster_cloud_logging_enabled" {
  type = bool
  default = true
}

variable "cluster_kubernetes_version" {
  type = string
  default = "v1.31.8"
}

variable "cluster_name" {
  type = string
  default = "cluster-name"
}

variable "cluster_security_group_id_list" {
  type = list(string)
  default = ["8a463aa4-b1dc-4f27-9c3f-53b94dc45e74","8a463aa4-b1dc-4f27-9c3f-53b94dc45e75"]
}

variable "cluster_subnet_id" {
  type = string
  default = "8a463aa4b1dc4f279c3f53b94dc45e74"
}

variable "cluster_vpc_id" {
  type = string
  default = "8a463aa4b1dc4f279c3f53b94dc45e74"
}

variable "cluster_volume_id" {
  type = string
  default = "8a463aa4-b1dc-4f27-9c3f-53b94dc45e74"
}

variable "private_endpoint_access_control_resources" {
  type = list(object({
    id = string
    name = string
    type = string
  }))
  default = [
    {
      id = "8a463aa4-b1dc-4f27-9c3f-53b94dc45e74"
      name = "name"
      type = "vm"
    }
  ]
}

variable "public_endpoint_access_control_ip" {
  type = string
  default = "10.10.10.10"
}