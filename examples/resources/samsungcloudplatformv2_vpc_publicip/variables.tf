variable "publicip_description" {
  type    = string
  default = "description-by-terrafrom-1211"
}

variable "publicip_type" {
  type    = string
  default = "GGW"
}

variable "zone" {
  type    = string
  default = "kr-west1-a"
}

variable "tags" {
  type = map(string)
  default = {
    tags1 = "testing"
  }
}




