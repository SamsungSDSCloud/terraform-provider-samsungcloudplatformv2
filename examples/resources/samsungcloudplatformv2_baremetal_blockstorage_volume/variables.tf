variable "region" {
  type    = string
  default = "kr-west1"
}

variable "block_storage_name" {
  type    = string
  default = "terraform-bs-01"
}

variable "disk_type" {
  type    = string
  default = "SSD"
}

variable "size_gb" {
  type    = number
  default = 10
}

variable "attachments" {
  type = list(object({
    object_type = string
    object_id   = string
  }))
  default = []
}

variable "qos" {
  type = object({
    iops       = number
    throughput = number
  })
  default = {
    iops       = 5000
    throughput = 250
  }
}

variable "tags" {
  type = map(string)
  default = {
    no_value = ""
    tf_key   = "tf_value"
  }
}

variable "create_timeouts" {
  type    = string
  default = "20m"
}

variable "delete_timeouts" {
  type    = string
  default = "20m"
}




