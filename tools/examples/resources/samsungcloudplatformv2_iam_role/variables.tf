variable "name" {
  type    = string
  default = "terraform-role"
}

variable "description" {
  type    = string
  default = "change test description"
}

variable "max_session_duration" {
  type    = number
  default = 3600
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
      action    = ["sts:AssumeRole"]
      condition = {}
      effect    = "Allow"
      principal = {
        principal_map = {
          Account = ["afd580f490394896a6bceabf77683c6b"]
        }
      }
      resource = ["*"]
      sid      = "VisualEditor0"
    }]
    version = "2024-07-01"
  }
}

variable "principals" {
  type = list(object({
    type  = string,
    value = string
  }))
  default = null
}

variable "policy_ids" {
  type    = list(string)
  default = ["ENTER YOUR RESOURCE'S POLICY_IDS"]
}

variable "tags" {
  type = map(string)
  default = {
    test_terraform_tag_key = "test_terraform_tag_value"
  }
}



