output "response" {
  value = {
    count : data.samsungcloudplatformv2_gslb_gslb_rrcs.my_gslb_rrc_list.total_count
    page : data.samsungcloudplatformv2_gslb_gslb_rrcs.my_gslb_rrc_list.page
    size : data.samsungcloudplatformv2_gslb_gslb_rrcs.my_gslb_rrc_list.size
    regional_gslbs : data.samsungcloudplatformv2_gslb_gslb_rrcs.my_gslb_rrc_list.regional_gslbs
  }
}