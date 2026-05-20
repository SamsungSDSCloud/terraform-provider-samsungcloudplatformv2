output "subnet_vips" {
  value = {
    count : data.samsungcloudplatformv2_vpc_subnet_vips.subnet_vip_list.total_count
    page : data.samsungcloudplatformv2_vpc_subnet_vips.subnet_vip_list.page
    size : data.samsungcloudplatformv2_vpc_subnet_vips.subnet_vip_list.size
    sort : data.samsungcloudplatformv2_vpc_subnet_vips.subnet_vip_list.sort
    subnet_vips : data.samsungcloudplatformv2_vpc_subnet_vips.subnet_vip_list.subnet_vips
  }
}
