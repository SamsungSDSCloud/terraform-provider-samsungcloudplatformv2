variable "name" {
  type = string
  default = "terraform-sg"
}
variable "description" {
  type = string
  default = "description info"
}

variable "loggable" {
  type = bool
  default = false
}

variable "security_group_tags" {
  type    = map(string)
  default = {
    "tf_key1" = "tf_val1",
    "tf_key2" = "tf_val2"
  }
}