variable "loadbalancer_id" {
  type    = string
  default = "ENTER YOUR RESOURCE'S LOADBALANCER_ID"
}
variable "static_nat" {
  type = object({
    publicip_id = string
  })
  default = {
    publicip_id = "ENTER YOUR RESOURCE'S PUBLICIP_ID"
  }
}




