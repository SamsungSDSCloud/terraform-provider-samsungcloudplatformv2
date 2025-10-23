variable "name" {
  type    = string
  default = ""
}

variable "volumes_filter_name" {
  type    = string
  default = ""
}

variable "volumes_filter_values" {
  type    = list(string)
  default = [""]
}

variable "volumes_filter_use_regex" {
  type    = bool
  default = false
}

