variable "loadbalancer_id" {
  type    = string
  default = ""
}
variable "static_nat" {
  type = object({
    private_nat_id    = string
    private_nat_ip_id = string
  })
  default = {
    private_nat_id    = ""
    private_nat_ip_id = ""
  }
}



