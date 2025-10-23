variable "account_id" {
  type    = string
  default = ""
}

variable "user_name" {
  type    = string
  default = ""
}

variable "description" {
  type    = string
  default = ""
}

variable "password" {
  type    = string
  default = ""
}

variable "temporary_password" {
  type    = bool
  default = false
}

variable "tags" {
  type    = map(string)
  default = null
}

variable "group_ids" {
  type    = list(string)
  default = [""]
}

variable "policy_ids" {
  type    = list(string)
  default = [""]
}

variable "password_reuse_count" {
  type    = number
  default = 0
}

