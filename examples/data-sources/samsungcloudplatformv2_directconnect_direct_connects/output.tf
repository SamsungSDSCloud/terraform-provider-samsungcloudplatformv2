output "directConnects" {
  value = {
    count: data.samsungcloudplatformv2_directconnect_direct_connects.directconnects.total_count,
    direct_connects: data.samsungcloudplatformv2_directconnect_direct_connects.directconnects.direct_connects,
    page: data.samsungcloudplatformv2_directconnect_direct_connects.directconnects.page,
    size: data.samsungcloudplatformv2_directconnect_direct_connects.directconnects.size,
    sort: data.samsungcloudplatformv2_directconnect_direct_connects.directconnects.sort_final,
  }
}