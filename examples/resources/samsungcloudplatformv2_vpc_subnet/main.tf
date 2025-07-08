provider "samsungcloudplatformv2" {
}

resource "samsungcloudplatformv2_vpc_subnet" "subnet" {
    name = var.subnet_name
    vpc_id = var.vpc_id
    type = var.subnet_type
    cidr = var.subnet_cidr
    description = var.subnet_description
    allocation_pools = var.subnet_allocation_pools
    dns_nameservers = var.subnet_dns_nameservers
    host_routes = var.subnet_host_routes
}