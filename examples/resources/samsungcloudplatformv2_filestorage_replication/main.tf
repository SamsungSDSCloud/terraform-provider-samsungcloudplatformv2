provider "samsungcloudplatformv2" {
}

resource "samsungcloudplatformv2_filestorage_replication" "replication" {
  backup_retention_count = var.backup_retention_count
  cifs_password          = var.cifs_password
  name                   = var.name
  region                 = var.region
  replication_frequency  = var.replication_frequency
  replication_type       = var.replication_type
  volume_id              = var.volume_id
  replication_update_type = var.replication_update_type
  replication_policy = var.replication_policy
}