provider "samsungcloudplatformv2" {
}

resource "samsungcloudplatformv2_vpc_port" "port" {
    name = var.port_name
    subnet_id = var.subnet_id
    description = var.port_description
    fixed_ip_address = var.port_fixed_ip_address
    security_groups = var.port_security_groups
}
