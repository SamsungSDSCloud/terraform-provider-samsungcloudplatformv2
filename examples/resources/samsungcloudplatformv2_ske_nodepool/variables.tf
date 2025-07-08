variable "name" {
  type    = string
  default = ""
}

variable "cluster_id" {
  type    = string
  default = ""
}

variable "desired_node_count" {
  type    = number
  default = 0
}

variable "image_os" {
  type    = string
  default = ""
}

variable "image_os_version" {
  type    = string
  default = ""
}

variable "is_auto_recovery" {
  type    = bool
  default = false
}

variable "is_auto_scale" {
  type    = bool
  default = false
}

variable "keypair_name" {
  type    = string
  default = ""
}

variable "kubernetes_version" {
  type    = string
  default = ""
}

variable "max_node_count" {
  type    = number
  default = 0
}

variable "min_node_count" {
  type    = number
  default = 0
}

variable "server_type_id" {
  type    = string
  default = ""
}

variable "volume_type_name" {
  type    = string
  default = ""
}

variable "volume_size" {
  type    = number
  default = 0
}

variable "labels" {
  type = list(object({
    key   = string
    value = string
  }))
  default = [{
    key   = ""
    value = ""
  }]
}

variable "taints" {
  type = list(object({
    effect = string
    key    = string
    value  = string
  }))
  default = [{
    effect = ""
    key    = ""
    value  = ""
  }]
}

