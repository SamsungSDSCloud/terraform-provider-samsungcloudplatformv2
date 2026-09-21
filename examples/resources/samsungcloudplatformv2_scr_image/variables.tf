variable "id" {
  type        = string
  default     = "ENTER YOUR RESOURCE'S ID"
  description = "The ID of the Docker image."
}

variable "description" {
  type        = string
  default     = "updated"
  description = "A description of the image."
}

variable "pull_policy" {
  type = object({
    critical_limit                  = number
    high_limit                      = number
    unmodified_excepted             = bool
    unscanned_image_pull_prevented  = bool
    vulnerable_image_pull_prevented = bool
  })
  default = {
    critical_limit                  = null
    high_limit                      = null
    unmodified_excepted             = false
    unscanned_image_pull_prevented  = false
    vulnerable_image_pull_prevented = true
  }
  description = "Pull policy configuration."
}

variable "scan_policy" {
  type = object({
    auto_scan_enabled      = bool
    fixed_version_excepted = bool
    language_excepted      = bool
    scan_policy_enabled    = bool
    secret_excepted        = bool
    severity_limit         = string
  })
  default = {
    auto_scan_enabled      = true
    fixed_version_excepted = false
    language_excepted      = true
    scan_policy_enabled    = true
    secret_excepted        = true
    severity_limit         = "None"
  }
  description = "Scan policy configuration."
}

variable "lifecycle_policy" {
  type = object({
    lifecycle_policy_enabled     = bool
    outdated_rule_duration       = number
    outdated_rule_enabled        = bool
    outdated_rule_tag_expression = string
    untagged_rule_duration       = number
    untagged_rule_enabled        = bool
  })
  default = {
    lifecycle_policy_enabled     = true
    outdated_rule_duration       = 360
    outdated_rule_enabled        = true
    outdated_rule_tag_expression = "*"
    untagged_rule_duration       = 360
    untagged_rule_enabled        = true
  }
  description = "Lifecycle policy configuration."
}

variable "lock_policy" {
  type = object({
    locked = bool
  })
  default = {
    locked = true
  }
  description = "Lock policy configuration."
}



