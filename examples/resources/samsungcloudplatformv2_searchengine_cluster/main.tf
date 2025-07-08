provider "samsungcloudplatformv2" {
}

resource "samsungcloudplatformv2_searchengine_cluster" "cluster" {
  allowable_ip_addresses  = var.allowable_ip_addresses
  dbaas_engine_version_id = var.dbaas_engine_version_id
  init_config_option      = var.init_config_option
  instance_groups         = var.instance_groups
  instance_name_prefix    = var.instance_name_prefix
  is_combined             = var.is_combined
  nat_enabled             = var.nat_enabled
  license                 = var.license
  name                    = var.name
  subnet_id               = var.subnet_id
  tags                    = var.tags
  service_state           = var.service_state
  timezone                = var.timezone
  maintenance_option      = var.maintenance_option
}