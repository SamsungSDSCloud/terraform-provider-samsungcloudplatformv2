variable "name" {
  type    = string
  default = ""
}

variable "size" {
  type    = number
  default = 0
}

variable "volume_type" {
  type    = string
  default = ""
}

variable "volume_server" {
  type = list(object({
    id = string
  }))
  default = [{
    id = ""
  }]
}

