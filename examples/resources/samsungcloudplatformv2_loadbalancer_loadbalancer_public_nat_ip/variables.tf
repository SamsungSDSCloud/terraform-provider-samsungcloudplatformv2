variable "loadbalancer_id" {
  type    = string
  default = ""
}
variable "static_nat" {
  type = object({
    publicip_id = string
  })
  default = {
    publicip_id = ""
  }
}



