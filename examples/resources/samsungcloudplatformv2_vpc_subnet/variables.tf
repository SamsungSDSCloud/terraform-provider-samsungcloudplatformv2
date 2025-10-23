variable "subnet_name" {
  type    = string
  default = ""
}

variable "vpc_id" {
  type    = string
  default = ""
}

variable "subnet_type" {
  type    = string
  default = ""
}

variable "subnet_cidr" {
  type    = string
  default = ""
}

variable "subnet_description" {
  type    = string
  default = ""
}

variable "subnet_allocation_pools" {
  type = list(object({
    start = string
    end   = string
  }))
  default = [{
    end   = ""
    start = ""
  }]
}

variable "subnet_dns_nameservers" {
  type    = list(string)
  default = [""]
}

variable "subnet_host_routes" {
  type = list(object({
    destination = string
    nexthop     = string
  }))
  default = [{
    destination = ""
    nexthop     = ""
  }]
}

