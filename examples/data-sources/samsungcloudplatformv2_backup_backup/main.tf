provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_backup_backup" "backup" {
  region = var.region

  id = var.id
  server_name = var.server_name
  name = var.name

  filter {
    name = var.backup_filter_name
    values = var.backup_filter_values
    use_regex = var.backup_filter_use_regex
  }
}
