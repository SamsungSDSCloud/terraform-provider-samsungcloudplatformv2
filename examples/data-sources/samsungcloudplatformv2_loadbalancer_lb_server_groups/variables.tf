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

variable "vpc_id" {
  type = string
  default = ""
}

variable "subnet_id" {
  type = string
  default = ""
}

# variable "protocol" {
#   type    = list(string)
#   default = ["TCP", "UDP"]
# }
