variable "budget_filter_name" {
  type    = string
  default = "name"
}

variable "budget_filter_values" {
  type    = list(string)
  default = ["budget_test"]
}

variable "budget_filter_use_regex" {
  type    = bool
  default = false
}


