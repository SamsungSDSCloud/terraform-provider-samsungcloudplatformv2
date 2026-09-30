provider "samsungcloudplatformv2" {
}

resource "samsungcloudplatformv2_smp_notification_group" "notification_group" {
  name               = var.notification_group_name
  description        = var.description
  recipient_user_ids = var.recipient_user_ids
}
