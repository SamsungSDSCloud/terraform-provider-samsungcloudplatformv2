variable "firewall_id" {
  type    = string
  default = "ENTER YOUR RESOURCE'S FIREWALL_ID"
}

variable "firewall_rule1" {
  type = object({
    action              = string
    description         = string
    destination_address = list(string)
    direction           = string
    service = list(object({
      service_type  = string
      service_value = string
    }))
    source_address = list(string)
    status         = string
  })
  default = {
    action              = "ALLOW"
    description         = "firewall rule 1"
    destination_address = ["1.1.1.1", "2.2.2.4"]
    direction           = "INBOUND"
    service = [{
      service_type  = "TCP"
      service_value = "443"
    }]
    source_address = ["4.4.4.4-4.4.4.10", "5.5.5.0/24"]
    status         = "ENABLE"
  }
}

variable "firewall_rule2" {
  type = object({
    action              = string
    description         = string
    destination_address = list(string)
    direction           = string
    service = list(object({
      service_type  = string
      service_value = string
    }))
    source_address = list(string)
    status         = string
  })
  default = {
    action              = "ALLOW"
    description         = "firewall rule 2"
    destination_address = ["4.4.4.4-4.4.4.10"]
    direction           = "OUTBOUND"
    service = [{
      service_type  = "TCP"
      service_value = "80"
    }]
    source_address = ["1.1.1.1", "2.2.2.2"]
    status         = "ENABLE"
  }
}

variable "firewall_rule3" {
  type = object({
    action              = string
    description         = string
    destination_address = list(string)
    direction           = string
    service = list(object({
      service_type  = string
      service_value = string
    }))
    source_address = list(string)
    status         = string
  })
  default = {
    action              = "DENY"
    description         = "firewall rule 3"
    destination_address = ["4.4.4.0/24"]
    direction           = "OUTBOUND"
    service = [{
      service_type  = "TCP"
      service_value = "80"
    }]
    source_address = ["1.1.1.1", "2.2.2.2"]
    status         = "ENABLE"
  }
}



