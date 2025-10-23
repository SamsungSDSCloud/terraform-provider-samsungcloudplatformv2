variable "name" {
  type    = string
  default = ""
}

variable "description" {
  type    = string
  default = ""
}

variable "max_session_duration" {
  type    = number
  default = 0
}

variable "assume_role_policy_document" {
  type = object({
    version = optional(string)
    statement = list(object({
      sid        = string
      effect     = string
      resource   = list(string)
      action     = optional(list(string))
      not_action = optional(list(string))
      condition  = optional(any)
      principal = optional(object({
        principal_string = optional(string)
        principal_map    = optional(map(list(string)))
      }))
    }))
  })
  nullable = true
  default = {
    statement = [{
      action     = null
      condition  = null
      effect     = ""
      not_action = null
      principal  = null
      resource   = [""]
      sid        = ""
    }]
    version = null
  }
}

variable "principals" {
  type = list(object({
    type  = string,
    value = string
  }))
  default = [{
    type  = ""
    value = ""
  }]
}

variable "policy_ids" {
  type    = list(string)
  default = [""]
}

variable "tags" {
  type    = map(string)
  default = null
}


