variable "resource_group_tags" {
  type    = map(string)
  default = {
    "tf_key1" = "tf_val1"
  }
}

variable "resource_group_filter_name" {
  type    = string
  default = ""
}

variable "resource_group_filter_values" {
  type    = list(string)
  default = [""]
}

variable "resource_group_filter_use_regex" {
  type    = bool
  default = true
}

