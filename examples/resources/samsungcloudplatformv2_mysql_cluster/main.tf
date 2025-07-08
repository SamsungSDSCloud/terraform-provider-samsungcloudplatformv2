provider "samsungcloudplatformv2" {
}

resource "samsungcloudplatformv2_mysql_cluster" "cluster" {
  dbaas_engine_version_id   = var.dbaas_engine_version_id
  allowable_ip_addresses    = var.allowable_ip_addresses
  ha_enabled                = var.ha_enabled
  init_config_option        = var.init_config_option
  instance_groups           = var.instance_groups
  instance_name_prefix      = var.instance_name_prefix
  name                      = var.name
  nat_enabled               = var.nat_enabled
  subnet_id                 = var.subnet_id
  tags                      = var.tags
  timezone                  = var.timezone
  maintenance_option        = var.maintenance_option
  service_state             = var.service_state
  virtual_ip_address        = var.virtual_ip_address
  vip_public_ip_id          = var.vip_public_ip_id
}