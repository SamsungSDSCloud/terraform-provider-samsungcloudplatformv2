variable "name" {
  type    = string
  default = null
}
variable "id" {
  type    = string
  default = "ENTER YOUR RESOURCE'S ID"
}
variable "volumes_filter_name" {
  type    = string
  default = "name"
}

variable "volumes_filter_values" {
  type    = list(string)
  default = ["test"]
}

variable "volumes_filter_use_regex" {
  type    = bool
  default = true
}


