variable "loadbalancer_id"{
  type = string
  default = "8a463aa4b1dc4f279c3f53b94dc45e74"
}
variable "static_nat" {
  type = object({
    publicip_id= string
  })
  default = {
    publicip_id= "8a463aa4b1dc4f279c3f53b94dc45e74"
    }
}

