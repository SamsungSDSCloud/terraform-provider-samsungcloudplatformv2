variable "policy_name" {
  type    = string
  default = "terraform-policy"
}

variable "policy_description" {
  type    = string
  default = "test description"
}

variable "policy_tags" {
  type = map(string)
  default = {
    test_terraform_tag_key = "test_terraform_tag_value"
  }
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
  default = {
    policy_document = {
      statement = [{
        action   = ["*"]
        effect   = "Allow"
        resource = ["*"]
        sid      = "VisualEditor0"
      }]
      version = "2024-07-01"
    }
  }
}



