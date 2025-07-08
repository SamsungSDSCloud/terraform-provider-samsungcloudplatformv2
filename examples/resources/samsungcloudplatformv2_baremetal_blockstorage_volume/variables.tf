variable "region" {
  type    = string
  default = ""
}

variable "block_storage_name" {
  type    = string
  default = ""
}

variable "disk_type" {
  type    = string
  default = ""
}

variable "size_gb" {
  type    = number
  default = 0
}

variable "attachments" {
  type = list(object({
    object_type = string
    object_id   = string
  }))
  default = [{
    object_id   = ""
    object_type = ""
  }]
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



