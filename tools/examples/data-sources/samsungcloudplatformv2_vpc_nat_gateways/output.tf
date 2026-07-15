output "natgateways" {
  value = {
    count : data.samsungcloudplatformv2_vpc_nat_gateways.natgateways.total_count,
    nat_gateways : data.samsungcloudplatformv2_vpc_nat_gateways.natgateways.nat_gateways,
    page : data.samsungcloudplatformv2_vpc_nat_gateways.natgateways.page,
    size : data.samsungcloudplatformv2_vpc_nat_gateways.natgateways.size,
    sort : data.samsungcloudplatformv2_vpc_nat_gateways.natgateways.sort_final,
  }
}