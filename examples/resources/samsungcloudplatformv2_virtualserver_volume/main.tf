provider "samsungcloudplatformv2" {
}

resource "samsungcloudplatformv2_virtualserver_volume" "volume" {
  name = var.name
  size = var.size
  volume_type = var.volume_type
  servers = var.volume_server

  tags = {
    "test_terraform_tag_key": "test_terraform_tag_value"
  }
}