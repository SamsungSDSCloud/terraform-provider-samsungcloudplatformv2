provider "samsungcloudplatformv2" {
}

resource "samsungcloudplatformv2_virtualserver_server" "server" {
  name            = var.name
  state           = var.state
  image_id        = var.image_id
  server_type_id  = var.server_type_id
  keypair_name    = var.keypair_name
  lock            = var.lock
  user_data       = var.user_data
  boot_volume     = var.boot_volume
  server_group_id = var.server_group_id
  security_groups = var.security_groups

  networks        = {
    interface_1 : {
      subnet_id : var.networks_interface_1_subnet_id,
    },
  }

  extra_volumes = {
    volume_1 : {
      size = var.extra_volumes_volume_1_size,
      type = var.extra_volumes_volume_1_type,
      delete_on_termination: var.extra_volumes_volume_1_delete_on_termination
    },
  }

  tags = {
    "test_terraform_tag_key": "test_terraform_tag_value"
  }
}
