output "vpcs" {
  value = {
    count: data.samsungcloudplatformv2_vpc_vpcs.vpcs.total_count
    page: data.samsungcloudplatformv2_vpc_vpcs.vpcs.page
    size: data.samsungcloudplatformv2_vpc_vpcs.vpcs.size
    sort: data.samsungcloudplatformv2_vpc_vpcs.vpcs.sort
    vpcs: data.samsungcloudplatformv2_vpc_vpcs.vpcs.vpcs
  }
}