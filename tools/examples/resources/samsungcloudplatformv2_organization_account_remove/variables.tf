variable "organization_id" {
  description = "Organization ID"
  type        = string
  default     = "ENTER YOUR RESOURCE'S ORGANIZATION_ID"
}

variable "account_id" {
  description = "Account ID to remove (single)"
  type        = string
  default     = "ENTER YOUR RESOURCE'S ACCOUNT_ID"
}

variable "target_account_ids" {
  description = "Account IDs to remove (multiple)"
  type        = list(string)
  default     = []
}


