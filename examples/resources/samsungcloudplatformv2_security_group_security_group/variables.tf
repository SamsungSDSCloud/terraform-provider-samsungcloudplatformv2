variable "name" {
  type    = string
  default = "sg-terraform-test-01"
}
variable "description" {
  type    = string
  default = "sg-terraform-test-01 description.."
}

variable "loggable" {
  type    = bool
  default = false
}

variable "security_group_tags" {
  type = map(string)
  default = {
    tf_key1 = "tf_val1"
    tf_key2 = "tf_val2"
  }
}


