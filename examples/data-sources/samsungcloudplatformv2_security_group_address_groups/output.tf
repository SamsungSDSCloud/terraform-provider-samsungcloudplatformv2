output "address_groups" {
  value = {
    size : data.samsungcloudplatformv2_security_group_address_groups.address_groups.size
    page : data.samsungcloudplatformv2_security_group_address_groups.address_groups.page
    count : data.samsungcloudplatformv2_security_group_address_groups.address_groups.total_count
    links : data.samsungcloudplatformv2_security_group_address_groups.address_groups.links
    address_groups : data.samsungcloudplatformv2_security_group_address_groups.address_groups.address_groups
  }
}
