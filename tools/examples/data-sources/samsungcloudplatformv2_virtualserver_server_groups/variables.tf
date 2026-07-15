variable "server_groups_filter_name" {
  type    = string
  default = "name"
}

variable "server_groups_filter_values" {
  type    = list(string)
  default = ["test"]
}

variable "server_groups_filter_use_regex" {
  type    = bool
  default = true
}


