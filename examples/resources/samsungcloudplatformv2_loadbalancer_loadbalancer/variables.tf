variable "loadbalancer" {
  type = object({
    description              = string
    firewall_enabled         = bool
    firewall_logging_enabled = bool
    layer_type               = string
    name                     = string
    service_ip               = string
    subnet_id                = string
    vpc_id                   = string
    source_nat_ip            = string
    health_check_ip_1        = string
    health_check_ip_2        = string
  })
  default = {
    description              = "remove public ip id param"
    firewall_enabled         = false
    firewall_logging_enabled = false
    health_check_ip_1        = null
    health_check_ip_2        = null
    layer_type               = "L4"
    name                     = "0908terraform-test-lb"
    service_ip               = null
    source_nat_ip            = null
    subnet_id                = "ENTER YOUR RESOURCE'S SUBNET_ID"
    vpc_id                   = "ENTER YOUR RESOURCE'S VPC_ID"
  }
}




