provider "samsungcloudplatformv2" {
}

resource "samsungcloudplatformv2_smp_rule_namespace" "rule_namespace" {
  name         = var.rule_namespace_name
  config_data  = var.config_data
  workspace_id = var.workspace_id
}
