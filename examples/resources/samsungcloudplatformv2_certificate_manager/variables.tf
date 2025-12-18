variable "cert_body" {
  type    = string
  default = ""
}

variable "cert_chain" {
  type    = string
  default = ""
}

variable "name" {
  type    = string
  default = ""
}

variable "private_key" {
  type    = string
  default = ""
}

variable "region" {
  type    = string
  default = ""
}

variable "timezone" {
  type    = string
  default = ""
}

variable "recipients" {
  type    = list(map(string))
  default = [null]
}

variable "tags" {
  type    = map(string)
  default = null
}

