variable "record" {
  type = object({
    description = string
    name        = string
    records     = list(string)
    ttl         = number
    type        = string
  })
  default = {
    description = "aaaaa"
    name        = "terraform22.com"
    records     = ["1.1.1.6"]
    ttl         = 300
    type        = "A"
  }
}

variable "id" {
  type    = string
  default = "ENTER YOUR RESOURCE'S ID"
}



