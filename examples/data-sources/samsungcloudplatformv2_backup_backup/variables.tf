variable "id" {
  type    = string
  default = "ENTER YOUR RESOURCE'S ID"
}

variable "server_name" {
  type    = string
  default = ""
}

variable "name" {
  type    = string
  default = ""
}

variable "backup_filter_name" {
  type    = string
  default = "name"
}

variable "backup_filter_values" {
  type    = list(string)
  default = ["backup"]
}

variable "backup_filter_use_regex" {
  type    = bool
  default = true
}


