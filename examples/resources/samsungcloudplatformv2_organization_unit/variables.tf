variable "name" {
  type        = string
  description = "Organization Unit name"
  default     = "test-org-unit"
}

variable "description" {
  type        = string
  description = "Organization Unit description"
  default     = "Test organization unit"
}

variable "organization_id" {
  type        = string
  description = "Organization ID"
  default     = "ENTER YOUR RESOURCE'S ORGANIZATION_ID"
}

variable "parent_unit_id" {
  type        = string
  description = "Parent Organization Unit ID (required)"
  default     = "ENTER YOUR RESOURCE'S PARENT_UNIT_ID"
}

variable "policy_ids" {
  type        = list(string)
  description = "Policy IDs"
  default     = ["ENTER YOUR RESOURCE'S POLICY_IDS"]
}


