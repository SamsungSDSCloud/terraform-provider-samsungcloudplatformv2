variable "loadbalancer" {
  type = object({
    description              = string
    firewall_enabled         = bool
    firewall_logging_enabled = bool
    layer_type               = string
    name                     = string
    publicip_id              = string
    service_ip               = string
    subnet_id                = string
    vpc_id                   = string
  })
  default = {
    description              = ""
    firewall_enabled         = false
    firewall_logging_enabled = false
    layer_type               = ""
    name                     = ""
    publicip_id              = ""
    service_ip               = ""
    subnet_id                = ""
    vpc_id                   = ""
  }
}



