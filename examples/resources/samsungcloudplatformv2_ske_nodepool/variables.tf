variable "name" {
  type    = string
  default = "terraform-test-2"
}

variable "cluster_id" {
  type    = string
  default = "ENTER YOUR RESOURCE'S CLUSTER_ID"
}

variable "desired_node_count" {
  type    = number
  default = 1
}

variable "image_os" {
  type    = string
  default = "ubuntu"
}

variable "image_os_version" {
  type    = string
  default = "22.04"
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
  default = "keypair-gook"
}

variable "kubernetes_version" {
  type    = string
  default = "v1.34.3"
}

variable "max_node_count" {
  type    = number
  default = null
}

variable "min_node_count" {
  type    = number
  default = null
}

variable "server_type_id" {
  type    = string
  default = "ENTER YOUR RESOURCE'S SERVER_TYPE_ID"
}

variable "volume_type_name" {
  type    = string
  default = "SSD"
}

variable "volume_size" {
  type    = number
  default = 104
}

variable "labels" {
  type = list(object({
    key   = string
    value = string
  }))
  default = [{
    key   = "test1"
    value = "test1"
    }, {
    key   = "test2"
    value = "test2"
    }, {
    key   = "test3"
    value = "test3"
    }, {
    key   = "test4"
    value = "test4"
  }]
}

variable "taints" {
  type = list(object({
    effect = string
    key    = string
    value  = string
  }))
  default = [{
    effect = "NoSchedule"
    key    = "test1"
    value  = "test1"
    }, {
    effect = "NoSchedule"
    key    = "test2"
    value  = "test2"
  }]
}

variable "server_group_id" {
  type    = string
  default = "ENTER YOUR RESOURCE'S SERVER_GROUP_ID"
}

variable "advanced_settings" {
  type = object({
    allowed_unsafe_sysctls  = string
    container_log_max_files = number
    container_log_max_size  = number
    image_gc_high_threshold = number
    image_gc_low_threshold  = number
    max_pods                = number
    pod_max_pids            = number
  })
  default = {
    allowed_unsafe_sysctls  = ""
    container_log_max_files = 5
    container_log_max_size  = 10
    image_gc_high_threshold = 85
    image_gc_low_threshold  = 80
    max_pods                = 110
    pod_max_pids            = 4096
  }
}


variable "linked_resources" {
  type = list(object({
    id   = string
    name = string
    type = string
  }))
  default = []
}

variable "volume_max_iops" {
  type    = number
  default = null
}

variable "volume_max_throughput" {
  type    = number
  default = null
}

variable "scp_gpu_driver" {
  type    = string
  default = null
}


