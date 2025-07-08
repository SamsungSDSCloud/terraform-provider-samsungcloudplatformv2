provider "samsungcloudplatformv2" {
}

resource "samsungcloudplatformv2_filestorage_volume" "volume" {
  name = var.name
  protocol = var.protocol
  type_name = var.type_name
  cifs_password = var.cifs_password
  file_unit_recovery_enabled = var.file_unit_recovery_enabled
  tags = var.tags
  access_rules = var.access_rules
}