variable "cluster_fabric_details" {
  type = object({
    cluster_fabric_id   = optional(string)
    cluster_fabric_name = string
    node_pool_id        = string
  })
  default = {
    cluster_fabric_id   = null
    cluster_fabric_name = ""
    node_pool_id        = ""
  }
}

variable "gpu_node_name_prefix" {
  type      = string
  ephemeral = true
  default   = ""
}

variable "image_id" {
  type    = string
  default = ""
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
  ephemeral = true
  default   = ""
}

variable "region_id" {
  type    = string
  default = ""
}

variable "server_type_id" {
  type      = string
  ephemeral = true
  default   = ""
}

variable "subnet_id" {
  type    = string
  default = ""
}


variable "server_details" {
  type = list(object({
    state = string
  }))
  default = [{
    state = ""
  }]
}

variable "vpc_id" {
  type    = string
  default = ""
}

variable "tags" {
  type    = map(string)
  default = null
}

variable "create_timeouts" {
  type    = string
  default = ""
}

variable "delete_timeouts" {
  type    = string
  default = ""
}


