variable "name" {
  type = string
  default = "terraform_volume"
}

variable "size" {
  type = number
  default = 8
}

variable "volume_type" {
  type = string
  default = "SSD"
}

variable "volume_server" {
  type = list(object({
    id = string
  }))
default = []
}