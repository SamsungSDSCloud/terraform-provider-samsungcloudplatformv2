variable "resource_groups_region" {
  type    = string
  default = "kr-west1"
}

variable "resource_groups_filter_name" {
  type    = string
  default = "name"
}

variable "resource_groups_filter_values" {
  type    = list(string)
  default = ["test"]
}

variable "resource_groups_filter_use_regex" {
  type    = bool
  default = true
}




