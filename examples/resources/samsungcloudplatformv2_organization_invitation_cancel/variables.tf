variable "organization_id" {
  description = "Organization ID"
  type        = string
  default     = "ENTER YOUR RESOURCE'S ORGANIZATION_ID"
}

variable "ids" {
  description = "Invitation IDs to cancel"
  type        = list(string)
  default     = ["ENTER YOUR RESOURCE'S IDS"]
}


