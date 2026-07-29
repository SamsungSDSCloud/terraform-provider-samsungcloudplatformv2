variable "name" {
  type    = string
  default = "terraform_test_volume"
}

variable "size" {
  type    = number
  default = 8
}

variable "volume_type" {
  type    = string
  default = "SSD"
}

variable "volume_server" {
  type = list(object({
    id = string
  }))
  default = []
}

variable "zone" {
  type    = string
  default = "kr-west1-a"
}


