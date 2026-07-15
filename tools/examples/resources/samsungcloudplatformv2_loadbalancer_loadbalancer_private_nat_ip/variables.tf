variable "loadbalancer_id" {
  type    = string
  default = "ENTER YOUR RESOURCE'S LOADBALANCER_ID"
}
variable "static_nat" {
  type = object({
    private_nat_id    = string
    private_nat_ip_id = string
  })
  default = {
    private_nat_id    = "ENTER YOUR RESOURCE'S PRIVATE_NAT_ID"
    private_nat_ip_id = "ENTER YOUR RESOURCE'S PRIVATE_NAT_IP_ID"
  }
}




