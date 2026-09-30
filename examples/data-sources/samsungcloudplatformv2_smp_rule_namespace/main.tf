provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_smp_rule_namespace" "rule_namespace" {
  rule_namespace_id = var.rule_namespace_id
  workspace_id      = var.workspace_id
}
