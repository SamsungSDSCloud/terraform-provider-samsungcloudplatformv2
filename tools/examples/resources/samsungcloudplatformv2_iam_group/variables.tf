variable "group_name" {
  type    = string
  default = "terraform-group"
}
variable "group_description" {
  type    = string
  default = "test description"
}
variable "group_tags" {
  type = map(string)
  default = {
    test_terraform_tag_key = "test_terraform_tag_value"
  }
}

variable "group_policy_ids" {
  type    = list(string)
  default = ["ENTER YOUR RESOURCE'S GROUP_POLICY_IDS"]
}

variable "group_user_ids" {
  type    = list(string)
  default = []
}


