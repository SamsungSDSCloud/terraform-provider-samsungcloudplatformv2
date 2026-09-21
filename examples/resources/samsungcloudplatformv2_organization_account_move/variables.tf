variable "organization_id" {
  description = "Organization ID"
  type        = string
  default     = "ENTER YOUR RESOURCE'S ORGANIZATION_ID"
}

variable "parent_unit_id" {
  description = "Target Parent Organization Unit ID"
  type        = string
  default     = "ENTER YOUR RESOURCE'S PARENT_UNIT_ID"
}

variable "target_account_ids" {
  description = "Account IDs to move"
  type        = list(string)
  default     = ["ENTER YOUR RESOURCE'S TARGET_ACCOUNT_IDS"]
}


