variable "private_dns" {
  type = object({
    connected_vpc_ids = list(string)
    description       = string
    name              = string
  })
  default = {
    connected_vpc_ids = "ENTER YOUR RESOURCE'S CONNECTED_VPC_IDS"
    description       = "aaaaa"
    name              = "terraform"
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


