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
    description              = ""
    firewall_enabled         = false
    firewall_logging_enabled = false
    health_check_ip_1        = ""
    health_check_ip_2        = ""
    layer_type               = ""
    name                     = ""
    service_ip               = ""
    source_nat_ip            = ""
    subnet_id                = ""
    vpc_id                   = ""
  }
}



