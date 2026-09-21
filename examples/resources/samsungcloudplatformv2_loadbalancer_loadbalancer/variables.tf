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
    health_check_ips         = list(string)
    zones                    = list(string)
  })
  default = {
    description              = "loadbalancer desc"
    firewall_enabled         = false
    firewall_logging_enabled = false
    health_check_ips         = null
    layer_type               = "L4"
    name                     = "nam-dep-trai-test-lb"
    service_ip               = null
    source_nat_ip            = null
    subnet_id                = "ENTER YOUR RESOURCE'S SUBNET_ID"
    vpc_id                   = "ENTER YOUR RESOURCE'S VPC_ID"
    zones                    = null
  }
}




