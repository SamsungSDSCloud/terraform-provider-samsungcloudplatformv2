provider "samsungcloudplatformv2" {
}

// Image From URL
resource "samsungcloudplatformv2_virtualserver_image" "image" {
  name              = var.name
  os_distro         = var.os_distro
  disk_format       = var.disk_format
  container_format  = var.container_format
  min_disk          = var.min_disk
  min_ram           = var.min_ram
  visibility        = var.visibility
  protected         = var.protected
  url               = var.url
  tags = {
    "test_terraform_tag_key": "test_terraform_tag_value"
  }
}

// Image From Server (Create)
resource "samsungcloudplatformv2_virtualserver_image" "image2" {
  name              = var.name
  instance_id       = var.instance_id
  tags = {}
}
