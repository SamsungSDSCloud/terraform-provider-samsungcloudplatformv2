variable "account_id" {
  description = "Organization Account ID"
  type        = string
}

variable "organization_id" {
  description = "Organization ID (optional query parameter)"
  type        = string
  default     = "ENTER YOUR RESOURCE'S ORGANIZATION_ID"
}

variable "lazy_policy" {
  description = "Linked Policy Query YN (optional query parameter)"
  type        = bool
  default     = null
}



