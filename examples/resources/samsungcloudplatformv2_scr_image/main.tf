provider "samsungcloudplatformv2" {
}

resource "samsungcloudplatformv2_scr_image" "image" {
  id          = var.id
  description = var.description

  pull_policy = var.pull_policy != null ? [
    {
      critical_limit                  = var.pull_policy.critical_limit
      high_limit                      = var.pull_policy.high_limit
      unmodified_excepted             = var.pull_policy.unmodified_excepted
      unscanned_image_pull_prevented  = var.pull_policy.unscanned_image_pull_prevented
      vulnerable_image_pull_prevented = var.pull_policy.vulnerable_image_pull_prevented
    }
  ] : []

  scan_policy = var.scan_policy != null ? [
    {
      auto_scan_enabled      = var.scan_policy.auto_scan_enabled
      fixed_version_excepted = var.scan_policy.fixed_version_excepted
      language_excepted      = var.scan_policy.language_excepted
      scan_policy_enabled    = var.scan_policy.scan_policy_enabled
      secret_excepted        = var.scan_policy.secret_excepted
      severity_limit         = var.scan_policy.severity_limit
    }
  ] : []

  lifecycle_policy = var.lifecycle_policy != null ? [
    {
      lifecycle_policy_enabled       = var.lifecycle_policy.lifecycle_policy_enabled
      outdated_rule_duration         = var.lifecycle_policy.outdated_rule_duration
      outdated_rule_enabled          = var.lifecycle_policy.outdated_rule_enabled
      outdated_rule_tag_expression   = var.lifecycle_policy.outdated_rule_tag_expression
      untagged_rule_duration         = var.lifecycle_policy.untagged_rule_duration
      untagged_rule_enabled          = var.lifecycle_policy.untagged_rule_enabled
    }
  ] : []

  lock_policy = var.lock_policy != null ? [
    {
      locked = var.lock_policy.locked
    }
  ] : []
}
