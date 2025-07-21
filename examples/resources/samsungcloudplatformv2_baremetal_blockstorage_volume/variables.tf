variable "region" {
  type = string
  default = "kr-west1"
}

variable "block_storage_name" {
  type = string
  default = "my-bs-01"
}

variable "disk_type" {
  type = string
  default = "SSD"
}

variable "size_gb" {
  type = number
  default = 10
}

variable "attachments" {
  type = list(object({
    object_type = string
    object_id = string
  }))
  default = [{
    object_type="BM",
    object_id="83c3c73d457345e3829ee6d5557c0011"
  }]
}

variable "tags" {
  type = map(string)
  default = {
    "tf_key": "tf_value",
    "no_value": ""
  }
}

variable "create_timeouts" {
  type = string
  default = "20m"
}

variable "delete_timeouts" {
  type = string
  default = "20m"
}

