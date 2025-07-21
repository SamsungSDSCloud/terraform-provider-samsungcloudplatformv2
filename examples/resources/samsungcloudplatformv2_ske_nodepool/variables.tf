variable "name" {
  type = string
  default = "np-name"
}

variable "cluster_id" {
  type = string
  default = "8a463aa4b1dc4f279c3f53b94dc45e74"
}

variable "desired_node_count" {
  type = number
  default = 2
}

variable "image_os" {
  type = string
  default = "ubuntu"
}

variable "image_os_version" {
  type = string
  default = "22.04"
}

variable "is_auto_recovery" {
  type = bool
  default = false
}

variable "is_auto_scale" {
  type = bool
  default = false
}

variable "keypair_name" {
  type = string
  default = "ssh"
}

variable "kubernetes_version" {
  type = string
  default = "v1.31.8"
}

variable "max_node_count" {
  type = number
  default = null
}

variable "min_node_count" {
  type = number
  default = null
}

variable "server_type_id" {
  type = string
  default = "s1v1m2"
}

variable "volume_type_name" {
  type = string
  default = "SSD"
}

variable "volume_size" {
  type = number
  default = 104
}

variable "labels" {
  type = list(object({
    key = string
    value = string
  }))
  default = [
    {
      key = "label1"
      value = "label1"
    },
    {
      key = "label2"
      value = "label2"
    }
  ]
}

variable "taints" {
  type = list(object({
    effect = string
    key = string
    value = string
  }))
  default = [
    {
      effect = "NoSchedule"
      key = "taint1"
      value = "taint1"
    },
    {
      effect = "NoSchedule"
      key = "taint2"
      value = "taint2"
    }
  ]
}