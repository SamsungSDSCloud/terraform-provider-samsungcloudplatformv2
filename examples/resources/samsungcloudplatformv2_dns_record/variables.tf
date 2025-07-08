variable "record" {
  type = object({
    description = string
    name        = string
    records     = list(string)
    ttl         = number
    type        = string
  })
  default = {
    description = ""
    name        = ""
    records     = [""]
    ttl         = 0
    type        = ""
  }
}

variable "id" {
  type    = string
  default = ""
}


