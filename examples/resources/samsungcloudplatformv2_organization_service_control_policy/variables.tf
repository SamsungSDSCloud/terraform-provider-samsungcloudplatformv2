variable "organization_id" {
  type        = string
  description = "Organization ID"
  default     = ""
}

variable "policy_name" {
  type        = string
  description = "Service Control Policy name"
  default     = ""
}

variable "policy_description" {
  type        = string
  description = "Service Control Policy description"
  default     = ""
}

variable "policy_type" {
  type        = string
  description = "Policy type (SYSTEM_MANAGED or USER_DEFINED)"
  default     = ""
}

variable "policy_document" {
  type = object({
    statement = list(object({
      effect     = string
      action     = list(string)
      not_action = list(string)
      principal  = string
      resource   = list(string)
      condition  = map(string)
      sid        = string
    }))
    version = string
  })
  description = "Service Control Policy document"
  default = {
    statement = [{
      action     = [""]
      condition  = null
      effect     = ""
      not_action = [""]
      principal  = ""
      resource   = [""]
      sid        = ""
    }]
    version = ""
  }
}

