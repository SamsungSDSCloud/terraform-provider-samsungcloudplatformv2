variable "group_name" {
  type    = string
  default = ""
}
variable "group_description" {
  type    = string
  default = ""
}
variable "group_tags" {
  type    = map(string)
  default = null
}

variable "group_policy_ids" {
  type    = list(string)
  default = [""]
}

variable "group_user_ids" {
  type    = list(string)
  default = [""]
}

