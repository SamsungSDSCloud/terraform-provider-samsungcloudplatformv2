variable "hosted_zone" {
  type = object({
    description    = string
    name           = string
    private_dns_id = string
    type           = string
  })
  default = {
    description    = ""
    name           = ""
    private_dns_id = ""
    type           = ""
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

