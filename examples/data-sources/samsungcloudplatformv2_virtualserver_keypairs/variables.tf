variable "keypairs_filter_name" {
  type    = string
  default = ""
}

variable "keypairs_filter_values" {
  type    = list(string)
  default = [""]
}

variable "keypairs_filter_use_regex" {
  type    = bool
  default = true
}