variable "organization_id" {
  description = "Organization ID"
  type        = string
  default     = ""
}

variable "account_id" {
  description = "Account ID to remove (single)"
  type        = string
  default     = ""
}

variable "target_account_ids" {
  description = "Account IDs to remove (multiple)"
  type        = list(string)
  default     = [""]
}

