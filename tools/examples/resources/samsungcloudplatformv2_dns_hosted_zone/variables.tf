variable "hosted_zone" {
  type = object({
    description    = string
    name           = string
    private_dns_id = string
    type           = string
  })
  default = {
    description    = "bbbbbbbbbbbbbb"
    name           = "terraform33.com"
    private_dns_id = "ENTER YOUR RESOURCE'S PRIVATE_DNS_ID"
    type           = "private"
  }
}

variable "tag" {
  type = object({
    test_terraform_tag_key = string
  })
  default = {
    test_terraform_tag_key = "test_terraform_tag_value"
  }
}


