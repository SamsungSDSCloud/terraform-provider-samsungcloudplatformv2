variable "id" {
  type    = string
  default = "ENTER YOUR RESOURCE'S ID"
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
  default = true
}

variable "size" {
  type    = number
  default = 10
}

variable "page" {
  type    = number
  default = 0
}

variable "sort" {
  type    = string
  default = "created_at:desc"
}

variable "state" {
  type    = list(string)
  default = []
}


