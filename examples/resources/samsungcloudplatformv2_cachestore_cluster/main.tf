provider "samsungcloudplatformv2" {
}

resource "samsungcloudplatformv2_cachestore_cluster" "cluster" {
  allowable_ip_addresses  = var.allowable_ip_addresses
  dbaas_engine_version_id = var.dbaas_engine_version_id
  ha_enabled              = var.ha_enabled
  init_config_option      = var.init_config_option
  instance_groups         = var.instance_groups
  instance_name_prefix    = var.instance_name_prefix
  maintenance_option      = var.maintenance_option
  name                    = var.name
  nat_enabled             = var.nat_enabled
  replica_count           = var.replica_count
  subnet_id               = var.subnet_id
  tags                    = var.tags
  timezone                = var.timezone
  service_state           = var.service_state
}