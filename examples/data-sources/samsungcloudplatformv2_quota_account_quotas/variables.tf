variable "account_quotas_filter_name" {
  type    = string
  default = "quota_item"
}

variable "account_quotas_filter_values" {
  type    = list(string)
  default = ["SECURITY_GROUP.ACCOUNT.RULE.MAX.COUNT"]
}

variable "account_quotas_filter_use_regex" {
  type    = bool
  default = false
}