variable "firewall_id" {
  type = string
  default = "8a463aa4b1dc4f279c3f53b94dc45e74"
}

variable "firewall_rule" {
  type = object({
    action = string
    description = string
    destination_address = list(string)
    direction = string
    service = list(object({
      service_type = string
      service_value = string
    }))
    source_address = list(string)
    status = string
  })
  default ={
    action = "ALLOW"
    description = "firewall description"
    destination_address = ["10.10.10.10","20.20.20.20"]
    direction = "INBOUND"
    service = [{
      service_type = "TCP"
      service_value = "443"
    },
    {
      service_type = "TCP"
      service_value = "3389"
    }]
    source_address = ["10.10.10.10-10.10.10.20","20.20.20.0/24"]
    status = "ENABLE"
  }
}
