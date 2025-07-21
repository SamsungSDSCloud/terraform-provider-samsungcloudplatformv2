variable "record" {
  type = object({
    description = string
    name        = string
    records = list(string)
    ttl         = number
    type        = string
  })
  default = {
    description = "description info"
    name        = "terraform.com"
    records = [
      "1.1.1.6"
    ]
    ttl  = 300
    type = "A"
  }
}

variable "id" {
  type = string
  default = "8a463aa4-b1dc-4f27-9c3f-53b94dc45e74"
}
