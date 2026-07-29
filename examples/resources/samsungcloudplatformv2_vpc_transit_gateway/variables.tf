variable "name" {
  default = "tgwNAMLEE01"
}

variable "tags" {
  type = map(string)
  default = {
    test_tag_key  = "test_tag_value"
    test_tag_key2 = "test_tag_value2"
  }
}

variable "description" {
  default = "description nam dep trai aaa"

}



