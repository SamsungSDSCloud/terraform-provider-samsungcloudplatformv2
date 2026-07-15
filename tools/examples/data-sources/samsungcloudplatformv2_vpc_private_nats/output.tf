output "privatenats" {
  value = {
    count : data.samsungcloudplatformv2_vpc_private_nats.privatenats.total_count
    page : data.samsungcloudplatformv2_vpc_private_nats.privatenats.page
    private_nats : data.samsungcloudplatformv2_vpc_private_nats.privatenats.private_nats
    size : data.samsungcloudplatformv2_vpc_private_nats.privatenats.size
    sort : data.samsungcloudplatformv2_vpc_private_nats.privatenats.sort
  }
}