variable "account_id" {
  type    = string
  default = "ENTER YOUR RESOURCE'S ACCOUNT_ID"
}

variable "user_name" {
  type    = string
  default = "terraform-user"
}

variable "description" {
  type    = string
  default = "description"
}

variable "password" {
  type    = string
  default = "ENTER YOUR RESOURCE'S PASSWORD"
}

variable "temporary_password" {
  type    = bool
  default = true
}

variable "tags" {
  type = map(string)
  default = {
    test_terraform_tag_key = "test_terraform_tag_value"
  }
}

variable "group_ids" {
  type    = list(string)
  default = []
}

variable "policy_ids" {
  type    = list(string)
  default = []
}

variable "password_reuse_count" {
  type    = number
  default = 2
}


