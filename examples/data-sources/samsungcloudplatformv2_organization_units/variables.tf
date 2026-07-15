variable "parent_unit_id" {
  type        = string
  description = "Parent Organization Unit ID (use 'root' for top-level)"
}

variable "organization_id" {
  type        = string
  description = "Organization ID"
  default     = "ENTER YOUR RESOURCE'S ORGANIZATION_ID"
}

variable "name" {
  type        = string
  description = "Organization Unit Name filter"
  default     = ""
}

variable "exclude_policy_id" {
  type        = string
  description = "Policy ID to exclude"
  default     = "ENTER YOUR RESOURCE'S EXCLUDE_POLICY_ID"
}


