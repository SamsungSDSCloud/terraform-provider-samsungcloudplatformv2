output "lbMembers" {
  value = {
    count : data.samsungcloudplatformv2_gslb_gslb_resources.gslbresources.total_count,
    gslb_resources : data.samsungcloudplatformv2_gslb_gslb_resources.gslbresources.gslb_resources,
    page : data.samsungcloudplatformv2_gslb_gslb_resources.gslbresources.page,
    size : data.samsungcloudplatformv2_gslb_gslb_resources.gslbresources.size,
    sort : data.samsungcloudplatformv2_gslb_gslb_resources.gslbresources.sort,
  }
}