variable "id" {
  type    = string
  default = "ENTER YOUR RESOURCE'S ID"
}

variable "server_group_filter_name" {
  type    = string
  default = "name"
}

variable "server_group_filter_values" {
  type    = list(string)
  default = ["test"]
}

variable "server_group_filter_use_regex" {
  type    = bool
  default = true
}


