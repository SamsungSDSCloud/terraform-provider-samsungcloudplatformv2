variable "invitation_id" {
  description = "Invitation ID to accept"
  type        = string
  default     = ""
}

variable "target_access_key" {
  description = "Target account access key"
  type        = string
  sensitive   = true
  default     = ""
}

variable "target_secret_key" {
  description = "Target account secret key"
  type        = string
  sensitive   = true
  default     = ""
}

