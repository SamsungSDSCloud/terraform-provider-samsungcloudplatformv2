variable "region" {
  type    = string
  default = "kr-west1"
}

variable "id" {
  type = string
  default = ""
}

variable "server_name" {
  type = string
  default = ""
}

variable "name" {
  type = string
  default = ""
}

variable "backup_filter_name" {
  type    = string
  default = ""
}

variable "backup_filter_values" {
  type    = list(string)
  default = ["backup"]
}

variable "backup_filter_use_regex" {
  type    = bool
  default = true
}