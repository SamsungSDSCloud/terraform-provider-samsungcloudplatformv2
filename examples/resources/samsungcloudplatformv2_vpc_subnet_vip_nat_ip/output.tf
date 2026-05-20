output "subnet_vip" {
  value = {
    id: samsungcloudplatformv2_vpc_subnet_vip_nat_ip.my_subnet_vip_nat_ip.id
    publicip_id: samsungcloudplatformv2_vpc_subnet_vip_nat_ip.my_subnet_vip_nat_ip.publicip_id
    state: samsungcloudplatformv2_vpc_subnet_vip_nat_ip.my_subnet_vip_nat_ip.state
  }
}
