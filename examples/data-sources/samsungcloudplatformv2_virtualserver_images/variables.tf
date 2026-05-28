variable "scp_image_type" {
  type    = string
  default = null
}

variable "scp_original_image_type" {
  type    = string
  default = null
}

variable "name" {
  type    = string
  default = null
}

variable "os_distro" {
  type    = string
  default = null
}

variable "status" {
  type    = string
  default = "active"
}

variable "visibility" {
  type    = string
  default = null
}

variable "images_filter_name" {
  type    = string
  default = "name"
}

variable "images_filter_values" {
  type    = list(string)
  default = ["test"]
}

variable "images_filter_use_regex" {
  type    = bool
  default = true
}


