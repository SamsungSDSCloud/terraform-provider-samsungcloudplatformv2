variable "region" {
  type    = string
  default = "kr-west1"
}

variable "server_name" {
  type = string
  default = ""
}

variable "name" {
  type = string
  default = ""
}

variable "backups_filter_name" {
  type    = string
  default = ""
}

variable "backups_filter_values" {
  type    = list(string)
  default = ["backup"]
}

variable "backups_filter_use_regex" {
  type    = bool
  default = true
}