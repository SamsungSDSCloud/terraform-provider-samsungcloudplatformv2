variable "id" {
  type    = string
  default = "ENTER YOUR RESOURCE'S ID"
}

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
  default = null
}

variable "visibility" {
  type    = string
  default = null
}

variable "image_filter_name" {
  type    = string
  default = "name"
}

variable "image_filter_values" {
  type    = list(string)
  default = ["test"]
}

variable "image_filter_use_regex" {
  type    = bool
  default = true
}


