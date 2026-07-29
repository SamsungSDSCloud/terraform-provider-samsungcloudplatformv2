variable "organization_id" {
  type        = string
  default     = "ENTER YOUR RESOURCE'S ORGANIZATION_ID"
  description = "Organization ID"
}

variable "document" {
  type = object({
    version = string
    statement = list(object({
      sid    = optional(string)
      effect = string
      action = list(string)
      principal = optional(object({
        scp = optional(list(string))
      }))
      resource = list(string)
    }))
  })
  default = {
    statement = [{
      action = ["organization:Create*", "organization:Delete*", "organization:Set*", "organization:Show*", "organization:List*", "organization:AttachPolicyBindings", "organization:RemovePolicyBindings"]
      effect = "Allow"
      principal = {
        scp = ["srn:qa2::a045159f40c64125a1fe61bd71d1c14a:::iam:user/b3a6e3f99c1040639e9a3c8f8b7427db"]
      }
      resource = ["*"]
      sid      = "statement1"
    }]
    version = "2012-10-17"
  }
  description = "Delegation policy document"
}


