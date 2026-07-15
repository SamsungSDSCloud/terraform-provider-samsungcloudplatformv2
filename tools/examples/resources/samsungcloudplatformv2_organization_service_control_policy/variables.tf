variable "organization_id" {
  type        = string
  default     = "ENTER YOUR RESOURCE'S ORGANIZATION_ID"
  description = "Organization ID"
}

variable "policy_name" {
  type        = string
  default     = "example-scp-policy"
  description = "Service Control Policy name"
}

variable "policy_description" {
  type        = string
  default     = "Example Service Control Policy"
  description = "Service Control Policy description"
}

variable "policy_type" {
  type        = string
  default     = "USER_DEFINED"
  description = "Policy type (SYSTEM_MANAGED or USER_DEFINED)"
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
  default = {
    statement = [{
      action     = ["samsungcloud:*"]
      condition  = {}
      effect     = "Allow"
      not_action = []
      principal  = "*"
      resource   = ["*"]
      sid        = "ExampleStatement"
    }]
    version = "2024-01-01"
  }
  description = "Service Control Policy document"
}


