variable "region" {
  type    = string
  default = "kr-west1"
}

variable "server_name" {
  type    = string
  default = "create-test"
}

variable "state" {
  type    = string
  default = "RUNNING"
}

variable "vpc_id" {
  type    = string
  default = "ENTER YOUR RESOURCE'S VPC_ID"
}

variable "ip" {
  type    = string
  default = ""
}

variable "baremetals_filter_name" {
  type    = string
  default = "server_name"
}

variable "baremetals_filter_values" {
  type    = list(string)
  default = ["test"]
}

variable "baremetals_filter_use_regex" {
  type    = bool
  default = false
}


