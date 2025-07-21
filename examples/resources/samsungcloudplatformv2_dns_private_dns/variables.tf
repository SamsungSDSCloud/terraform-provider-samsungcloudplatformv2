variable "private_dns" {
  type = object({
    connected_vpc_ids = list(string)
    description = string
    name        = string
  })
  default = {
    connected_vpc_ids = [
      "8a463aa4b1dc4f279c3f53b94dc45e74"
    ]
    description = "description info"
    name = "terraform"
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