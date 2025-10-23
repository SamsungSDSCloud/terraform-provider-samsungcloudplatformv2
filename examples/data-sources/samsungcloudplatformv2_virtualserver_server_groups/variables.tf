variable "server_groups_filter_name" {
  type    = string
  default = ""
}

variable "server_groups_filter_values" {
  type    = list(string)
  default = [""]
}

variable "server_groups_filter_use_regex" {
  type    = bool
  default = false
}

