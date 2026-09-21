provider "samsungcloudplatformv2" {
}

resource "samsungcloudplatformv2_scr_repository" "repository" {
  registry_id = var.registry_id
  name        = var.name
  description = var.description

  lifecycle_policy = [{
    lifecycle_policy_enabled    = true
    outdated_rule_enabled        = true
    outdated_rule_duration       = 90
    outdated_rule_tag_expression = "*"
    untagged_rule_enabled        = true
    untagged_rule_duration        = 7
  }]

  lock_policy = [{
    locked = true
  }]

  pull_policy = [{
    unmodified_excepted             = false
    unscanned_image_pull_prevented  = true
    vulnerable_image_pull_prevented = true
    critical_limit                  = 0
    high_limit                      = 10
  }]

  scan_policy = [{
    auto_scan_enabled      = true
    fixed_version_excepted = false
    language_excepted      = false
    scan_policy_enabled    = true
    secret_excepted        = false
    severity_limit         = "High"
  }]

  tags = {
    Environment = "production"
    Team        = "devops"
  }
}