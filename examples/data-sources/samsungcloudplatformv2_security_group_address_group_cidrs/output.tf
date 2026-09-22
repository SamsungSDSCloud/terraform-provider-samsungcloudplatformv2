output "address_groups" {
  value = {
    addresses: data.samsungcloudplatformv2_security_group_address_group_cidrs.addresses.addresses,
    count: data.samsungcloudplatformv2_security_group_address_group_cidrs.addresses.total_count,
    links: data.samsungcloudplatformv2_security_group_address_group_cidrs.addresses.links,
  }
}
