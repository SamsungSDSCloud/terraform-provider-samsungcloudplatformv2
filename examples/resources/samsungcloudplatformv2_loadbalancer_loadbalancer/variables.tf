variable "loadbalancer" {
  type = object({
    description = string
    firewall_enabled = bool
    firewall_logging_enabled= bool
    layer_type= string
    name= string
    publicip_id= string
    service_ip= string
    subnet_id= string
    vpc_id= string
  })
  default = {
    description = "description info"
    firewall_enabled = true
    firewall_logging_enabled= false
    layer_type= "L4"
    name= "terraform-lb"
    publicip_id= null
    service_ip= null
    subnet_id= "8a463aa4b1dc4f279c3f53b94dc45e74"
    vpc_id= "8a463aa4b1dc4f279c3f53b94dc45e74"
    }
}

