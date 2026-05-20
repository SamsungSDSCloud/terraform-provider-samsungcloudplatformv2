variable "organization_id" {
  description = "Organization ID"
  type        = string
  default     = ""
}

variable "ids" {
  description = "Invitation IDs to cancel"
  type        = list(string)
  default     = [""]
}

