provider "samsungcloudplatformv2" {
}

resource "samsungcloudplatformv2_baremetal_baremetal" "baremetal" {
  image_id = var.image_id
  init_script = var.init_script
  lock_enabled = var.lock_enabled
  os_user_id = var.os_user_id
  os_user_password = var.os_user_password
  placement_group_name = var.placement_group_name
  region_id = var.region_id
  server_details = var.server_details
  subnet_id = var.subnet_id
  use_placement_group = var.use_placement_group
  vpc_id = var.vpc_id
  tags = var.tags

  timeouts {
    create = var.create_timeouts
    delete = var.delete_timeouts
  }
}