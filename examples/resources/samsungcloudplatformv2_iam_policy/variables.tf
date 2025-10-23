variable "policy_name" {
  type    = string
  default = ""
}

variable "policy_description" {
  type    = string
  default = ""
}

variable "policy_tags" {
  type    = map(string)
  default = null
}

variable "policy_version" {
  type = map(object({
    version = string
    statement = list(object({
      sid        = string,
      effect     = string,
      resource   = list(string),
      action     = optional(list(string)),
      not_action = optional(list(string)),
      condition  = optional(any),
    }))
  }))
  default = null
}


