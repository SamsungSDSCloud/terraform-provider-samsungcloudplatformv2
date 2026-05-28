variable "cluster_fabric_details" {
  type = object({
    cluster_fabric_id   = optional(string)
    cluster_fabric_name = string
    node_pool_id        = string
  })
  default = {
    cluster_fabric_id   = "ENTER YOUR RESOURCE'S CLUSTER_FABRIC_ID"
    cluster_fabric_name = "clusterb300"
    node_pool_id        = "ENTER YOUR RESOURCE'S NODE_POOL_ID"
  }
}

variable "gpu_node_name_prefix" {
  type      = string
  ephemeral = true
  default   = "gpu-b300"
}

variable "image_id" {
  type    = string
  default = "ENTER YOUR RESOURCE'S IMAGE_ID"
}

variable "init_script" {
  type    = string
  default = ""
}

variable "lock_enabled" {
  type    = bool
  default = false
}

variable "os_user_password" {
  type      = string
  default   = "ENTER YOUR RESOURCE'S OS_USER_PASSWORD"
  ephemeral = true
}

variable "region_id" {
  type    = string
  default = "ENTER YOUR RESOURCE'S REGION_ID"
}

variable "server_type_id" {
  type      = string
  default   = "ENTER YOUR RESOURCE'S SERVER_TYPE_ID"
  ephemeral = true
}

variable "subnet_id" {
  type    = string
  default = "ENTER YOUR RESOURCE'S SUBNET_ID"
}


variable "server_details" {
  type = list(object({
    state = string
  }))
  default = [{
    state = "RUNNING"
  }]
}

variable "vpc_id" {
  type    = string
  default = "ENTER YOUR RESOURCE'S VPC_ID"
}

variable "tags" {
  type = map(string)
  default = {
    "1stopped" = "delsoon"
  }
}

variable "create_timeouts" {
  type    = string
  default = "60m"
}

variable "delete_timeouts" {
  type    = string
  default = "40m"
}



