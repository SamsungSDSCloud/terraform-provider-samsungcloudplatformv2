variable "hosted_zone" {
  type = object({
    description = string
    email       = string
    name        = string
    type        = string
  })
  default = {
    description = "description info"
    email       = "sung.sam@samsung.com"
    name        = "terraformhosted.com"
    type        = "public"
  }
}

variable "tag" {
  type = object({
    terraform_tag_key = string
  })
  default = {
    terraform_tag_key = "terraform_tag_value"
  }
}