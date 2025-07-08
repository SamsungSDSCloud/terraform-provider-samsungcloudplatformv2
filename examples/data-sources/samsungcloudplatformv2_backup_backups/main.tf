provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_backup_backups" "ids" {
  region = var.region

  server_name = var.server_name
  name = var.name

  filter {
    name = var.backups_filter_name
    values = var.backups_filter_values
    use_regex = var.backups_filter_use_regex
  }
}