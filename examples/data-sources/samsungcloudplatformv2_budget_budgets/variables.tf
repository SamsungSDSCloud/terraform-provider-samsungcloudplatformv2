variable "budget_filter_name" {
  type    = string
  default = ""
}

variable "budget_filter_values" {
  type    = list(string)
  default = [""]
}

variable "budget_filter_use_regex" {
  type    = bool
  default = false
}

