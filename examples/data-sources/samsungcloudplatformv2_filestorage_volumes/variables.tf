variable "limit" {
  type = number
  default = 10
}

variable "offset" {
  type = number
  default = 0
}

variable "sort" {
  type = string
  default = "created_at:desc"
}

variable "name" {
  type = string
  default = "test"
}

variable "type_name" {
  type = string
  default = "HDD"
}
