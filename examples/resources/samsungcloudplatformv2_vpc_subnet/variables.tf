variable "subnet_name" {
  type = string
  default = "subnetName"
}

variable "vpc_id" {
  type = string
  default = "7df8abb4912e4709b1cb237daccca7a8"
}

variable "subnet_type" {
  type = string
  default = "GENERAL"
}

variable "subnet_cidr" {
  type = string
  default = "192.167.1.0/24"
}

variable "subnet_description" {
  type = string
  default = "description info"
}

variable "subnet_allocation_pools" {
  type = list(object({
      start = string
      end = string
    }))
    default = [
      {
        start = "10.0.0.2"
        end = "10.0.0.254"
      }
    ]
}

variable "subnet_dns_nameservers" {
  type = list(string)
  default = ["1.1.1.1"]
}

variable "subnet_host_routes" {
  type = list(object({
      destination = string
      nexthop = string
    }))
    default = [
      {
        destination = "192.168.24.0/24"
        nexthop = "192.168.20.5"
      }
    ]
}