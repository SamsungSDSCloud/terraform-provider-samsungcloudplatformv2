variable "name" {
  type    = string
  default = ""
}
variable "description" {
  type    = string
  default = ""
}

variable "loggable" {
  type    = bool
  default = false
}

variable "security_group_tags" {
  type    = map(string)
  default = null
}

