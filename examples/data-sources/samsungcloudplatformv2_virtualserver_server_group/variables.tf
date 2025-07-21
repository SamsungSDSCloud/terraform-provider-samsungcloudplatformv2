variable "id" {
  type = string
  default = ""
}

variable "server_group_filter_name" {
  type    = string
  default = ""
}

variable "server_group_filter_values" {
  type    = list(string)
  default = [""]
}

variable "server_group_filter_use_regex" {
  type    = bool
  default = true
}