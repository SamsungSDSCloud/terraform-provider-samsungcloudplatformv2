variable "server_name" {
  type    = string
  default = ""
}

variable "name" {
  type    = string
  default = ""
}

variable "backups_filter_name" {
  type    = string
  default = "name"
}

variable "backups_filter_values" {
  type    = list(string)
  default = ["backup"]
}

variable "backups_filter_use_regex" {
  type    = bool
  default = true
}


