provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_smp_notification_groups" "notification_groups" {
  name       = var.name
  name_like  = var.name_like
  id         = var.id
  id_like    = var.id_like
  page       = var.page
  size       = var.size
  sort_by    = var.sort_by
  sort_order = var.sort_order
}
