variable "subnet_name" {
  type    = string
  default = "testsubnet"
}

variable "vpc_id" {
  type    = string
  default = "ENTER YOUR RESOURCE'S VPC_ID"
}

variable "subnet_type" {
  type    = string
  default = "GENERAL"
}

variable "subnet_cidr" {
  type    = string
  default = "192.168.0.0/28"
}

variable "subnet_description" {
  type    = string
  default = "test_description"
}

variable "subnet_allocation_pools" {
  type = list(object({
    start = string
    end   = string
  }))
  default = [{
    end   = "192.168.0.12"
    start = "192.168.0.10"
    }, {
    end   = "192.168.0.4"
    start = "192.168.0.3"
  }]
}

variable "subnet_dns_nameservers" {
  type    = set(string)
  default = ["8.8.8.8"]
}

variable "subnet_host_routes" {
  type = list(object({
    destination = string
    nexthop     = string
  }))
  default = [{
    destination = "192.168.24.0/24"
    nexthop     = "11.11.11.11"
    }, {
    destination = "192.169.24.0/24"
    nexthop     = "22.22.22.22"
  }]
}

variable "dhcp_ip_address" {
  type    = string
  default = null
}

variable "gateway_ip_address" {
  type    = string
  default = null
}

variable "tags" {
  type = map(string)
  default = {
    tf = "terraform"
  }
}


