variable "firewall_id" {
  type    = string
  default = ""
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
    action              = ""
    description         = ""
    destination_address = [""]
    direction           = ""
    service = [{
      service_type  = ""
      service_value = ""
    }]
    source_address = [""]
    status         = ""
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
    action              = ""
    description         = ""
    destination_address = [""]
    direction           = ""
    service = [{
      service_type  = ""
      service_value = ""
    }]
    source_address = [""]
    status         = ""
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
    action              = ""
    description         = ""
    destination_address = [""]
    direction           = ""
    service = [{
      service_type  = ""
      service_value = ""
    }]
    source_address = [""]
    status         = ""
  }
}


