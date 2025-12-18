variable "id" {
  type    = string
  default = ""
}

variable "name" {
  type    = string
  default = ""
}

variable "cn" {
  type    = string
  default = ""
}

variable "is_mine" {
  type    = bool
  default = false
}

variable "size" {
  type    = number
  default = 0
}

variable "page" {
  type    = number
  default = 0
}

variable "sort" {
  type    = string
  default = ""
}

variable "state" {
  type    = list(string)
  default = [""]
}

