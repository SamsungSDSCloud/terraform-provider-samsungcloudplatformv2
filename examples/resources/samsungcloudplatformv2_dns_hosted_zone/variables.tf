variable "hosted_zone" {
  type = object({
    description = string
    email       = string
    name        = string
    type        = string
  })
  default = {
    description = ""
    email       = ""
    name        = ""
    type        = ""
  }
}

variable "tag" {
  type = object({
    test_terraform_tag_key = string
  })
  default = {
    test_terraform_tag_key = ""
  }
}

