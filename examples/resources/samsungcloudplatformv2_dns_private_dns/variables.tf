variable "private_dns" {
  type = object({
    connected_vpc_ids = list(string)
    description       = string
    name              = string
  })
  default = {
    connected_vpc_ids = [""]
    description       = ""
    name              = ""
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

