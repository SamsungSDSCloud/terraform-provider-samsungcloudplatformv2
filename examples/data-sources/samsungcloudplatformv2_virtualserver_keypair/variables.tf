variable "name" {
  type = string
  default = "terraform_keypair"
}

variable "keypair_filter_name" {
  type    = string
  default = ""
}

variable "keypair_filter_values" {
  type    = list(string)
  default = [""]
}

variable "keypair_filter_use_regex" {
  type    = bool
  default = true
}