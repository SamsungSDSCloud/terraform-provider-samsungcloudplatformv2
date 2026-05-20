output "publicips" {
  value = {
    count: data.samsungcloudplatformv2_vpc_publicips.publicips.total_count,
    page: data.samsungcloudplatformv2_vpc_publicips.publicips.page,
    publicips: data.samsungcloudplatformv2_vpc_publicips.publicips.publicips,
    size: data.samsungcloudplatformv2_vpc_publicips.publicips.size,
    sort: data.samsungcloudplatformv2_vpc_publicips.publicips.sort,
  }
}
