provider "samsungcloudplatformv2" {
}

resource "samsungcloudplatformv2_vertica_cluster" "cluster" {
  name                    = var.name
  instance_name_prefix    = var.instance_name_prefix
  dbaas_engine_version_id = var.dbaas_engine_version_id
  subnet_id               = var.subnet_id
  instance_groups         = var.instance_groups
  allowable_ip_addresses  = var.allowable_ip_addresses
  init_config_option      = var.init_config_option
  maintenance_option      = var.maintenance_option
  service_state           = var.service_state
  nat_enabled             = var.nat_enabled
  timezone                = var.timezone
  license                 = var.license
  tags                    = var.tags
}