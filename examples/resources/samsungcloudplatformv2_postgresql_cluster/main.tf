provider "samsungcloudplatformv2" {
}

resource "samsungcloudplatformv2_postgresql_cluster" "cluster" {
  allowable_ip_addresses  = var.allowable_ip_addresses
  dbaas_engine_version_id = var.dbaas_engine_version_id
  nat_enabled             = var.nat_enabled
  ha_enabled              = var.ha_enabled
  init_config_option      = var.init_config_option
  instance_groups         = var.instance_groups
  instance_name_prefix    = var.instance_name_prefix
  name                    = var.name
  subnet_id               = var.subnet_id
  tags                    = var.tags
  service_state           = var.service_state
  timezone                = var.timezone
  maintenance_option      = var.maintenance_option
  vip_public_ip_id        = var.vip_public_ip_id
  virtual_ip_address      = var.virtual_ip_address
}