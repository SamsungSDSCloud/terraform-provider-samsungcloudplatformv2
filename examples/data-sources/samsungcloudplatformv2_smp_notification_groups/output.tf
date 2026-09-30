output "notification_groups" {
  value = data.samsungcloudplatformv2_smp_notification_groups.notification_groups
}

output "total_count" {
  value = data.samsungcloudplatformv2_smp_notification_groups.notification_groups.total_count
}

output "page" {
  value = data.samsungcloudplatformv2_smp_notification_groups.notification_groups.page
}

output "size" {
  value = data.samsungcloudplatformv2_smp_notification_groups.notification_groups.size
}

output "sort" {
  value = data.samsungcloudplatformv2_smp_notification_groups.notification_groups.sort
}
