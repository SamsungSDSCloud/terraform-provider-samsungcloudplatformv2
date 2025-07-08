variable "account_quotas_filter_name" {
  type    = string
  default = ""
}

variable "account_quotas_filter_values" {
  type    = list(string)
  default = [""]
}

variable "account_quotas_filter_use_regex" {
  type    = bool
  default = false
}

