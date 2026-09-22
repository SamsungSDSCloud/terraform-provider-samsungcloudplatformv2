variable "organization_id" {
  description = "Organization ID"
  type        = string
  default     = "ENTER YOUR RESOURCE'S ORGANIZATION_ID"
}

variable "target_login_ids" {
  description = "Target Login IDs to invite"
  type        = list(string)
  default     = ["ENTER YOUR RESOURCE'S TARGET_LOGIN_IDS"]
}


