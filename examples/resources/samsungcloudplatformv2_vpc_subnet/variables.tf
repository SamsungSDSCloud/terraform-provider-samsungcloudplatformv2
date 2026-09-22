variable "subnet_name" {
  type    = string
  default = "yb-subnet-2"
}

variable "vpc_id" {
  type    = string
  default = "ENTER YOUR RESOURCE'S VPC_ID"
}

variable "subnet_type" {
  type    = string
  default = "PUBLIC"
}

variable "subnet_cidr" {
  type    = string
  default = "30.30.16.0/21"
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
  default = null
}

variable "subnet_dns_nameservers" {
  type    = set(string)
  default = null
}

variable "subnet_host_routes" {
  type = list(object({
    destination = string
    nexthop     = string
  }))
  default = null
}

variable "dhcp_ip_address" {
  type    = string
  default = "30.30.16.3"
}

variable "gateway_ip_address" {
  type    = string
  default = null
}

variable "tags" {
  type    = map(string)
  default = null
}
variable "category" {
  type    = string
  default = "SECONDARY"
}
variable "primary_subnet_id" {
  type    = string
  default = "ENTER YOUR RESOURCE'S PRIMARY_SUBNET_ID"
}


