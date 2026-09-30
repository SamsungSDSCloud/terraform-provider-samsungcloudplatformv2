provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_smp_notification_group" "notification_group" {
  notification_group_id = var.notification_group_id
}
