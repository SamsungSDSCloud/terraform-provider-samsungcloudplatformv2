variable "name" {
  type    = string
  default = "terraform_test_keypair"
}

variable "keypair_filter_name" {
  type    = string
  default = "name"
}

variable "keypair_filter_values" {
  type    = list(string)
  default = ["test"]
}

variable "keypair_filter_use_regex" {
  type    = bool
  default = true
}


