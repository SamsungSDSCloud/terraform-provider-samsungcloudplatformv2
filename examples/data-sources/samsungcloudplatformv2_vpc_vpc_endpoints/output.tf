output "vpcendpoints" {
  value = {
    count : data.samsungcloudplatformv2_vpc_vpc_endpoints.vpcendpoints.total_count,
    page : data.samsungcloudplatformv2_vpc_vpc_endpoints.vpcendpoints.page,
    size : data.samsungcloudplatformv2_vpc_vpc_endpoints.vpcendpoints.size,
    sort : data.samsungcloudplatformv2_vpc_vpc_endpoints.vpcendpoints.sort,
    vpc_endpoints : data.samsungcloudplatformv2_vpc_vpc_endpoints.vpcendpoints.vpc_endpoints
  }
}
