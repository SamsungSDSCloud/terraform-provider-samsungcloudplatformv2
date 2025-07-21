variable "name" {
  type = string
  default = ""
}

variable "size" {
  type = number
  default = 20
}

variable "page" {
  type = number
  default = 0
}

variable "sort" {
  type = string
  default = ""
}

variable "subnet_id" {
  type = string
  default = "8a463aa4b1dc4f279c3f53b94dc45e74"
}

# variable "protocol" {
#   type    = list(string)
#   default = ["TCP", "UDP"]
# }
