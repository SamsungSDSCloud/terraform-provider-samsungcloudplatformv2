variable "resource_groups_region" {
  type    = string
  default = ""
}

variable "resource_groups_filter_name" {
  type    = string
  default = ""
}

variable "resource_groups_filter_values" {
  type    = list(string)
  default = [""]
}

variable "resource_groups_filter_use_regex" {
  type    = bool
  default = false
}



