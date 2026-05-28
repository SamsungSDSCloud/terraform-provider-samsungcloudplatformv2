variable "keypairs_filter_name" {
  type    = string
  default = "name"
}

variable "keypairs_filter_values" {
  type    = list(string)
  default = ["test"]
}

variable "keypairs_filter_use_regex" {
  type    = bool
  default = true
}


