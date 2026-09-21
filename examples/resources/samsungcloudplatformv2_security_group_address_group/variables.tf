variable "name" {
  type    = string
  default = "gdcv-terraform-test-02"
}
variable "description" {
  type    = string
  default = "123123"
}

variable "addresses" {
  type    = set(string)
  default = ["10.1.1.0/31"]
}

variable "tags" {
  type = map(string)
  default = {
    tf_key1 = "tf_val2"
  }
}


