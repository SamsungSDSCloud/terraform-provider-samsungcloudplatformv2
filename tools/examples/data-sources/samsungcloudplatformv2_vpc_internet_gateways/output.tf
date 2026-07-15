output "internetGateways" {
  value = {
    count : data.samsungcloudplatformv2_vpc_internet_gateways.internetgateways.total_count,
    internet_gateways : data.samsungcloudplatformv2_vpc_internet_gateways.internetgateways.internet_gateways,
    page : data.samsungcloudplatformv2_vpc_internet_gateways.internetgateways.page,
    size : data.samsungcloudplatformv2_vpc_internet_gateways.internetgateways.size,
    sort : data.samsungcloudplatformv2_vpc_internet_gateways.internetgateways.sort_final,
  }
}