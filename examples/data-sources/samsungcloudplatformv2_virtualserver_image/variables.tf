variable "id" {
  type = string
  default = ""
}

variable "scp_image_type" {
  type = string
  default = ""
}

variable "scp_original_image_type" {
  type = string
  default = ""
}

variable "name" {
  type = string
  default = ""
}

variable "os_distro" {
  type = string
  default = ""
}

variable "status" {
  type = string
  default = ""
}

variable "visibility" {
  type = string
  default = ""
}

variable "image_filter_name" {
  type    = string
  default = ""
}

variable "image_filter_values" {
  type    = list(string)
  default = [""]
}

variable "image_filter_use_regex" {
  type    = bool
  default = true
}