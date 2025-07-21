variable "size" {
  type = number
  default = 10
}

variable "page" {
  type = number
  default = 0
}

variable "sort" {
  type = string
  default = "created_at:desc"
}

variable "name" {
  type = string
  default = ""
}

variable "service_state" {
  type = string
  default = ""
}

variable "database_name" {
  type = string
  default = ""
}