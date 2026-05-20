variable "organization_id" {
  description = "Organization ID"
  type        = string
  default     = ""
}

variable "parent_unit_id" {
  description = "Target Parent Organization Unit ID"
  type        = string
  default     = ""
}

variable "target_account_ids" {
  description = "Account IDs to move"
  type        = list(string)
  default     = [""]
}

