variable "organization_id" {
  description = "Organization ID"
  type        = string
  default     = ""
}

variable "target_login_ids" {
  description = "Target Login IDs to invite"
  type        = list(string)
  default     = [""]
}

