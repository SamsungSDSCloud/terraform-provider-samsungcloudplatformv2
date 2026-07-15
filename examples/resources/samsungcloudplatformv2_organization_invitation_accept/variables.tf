variable "invitation_id" {
  description = "Invitation ID to accept"
  type        = string
}

variable "target_access_key" {
  description = "Target account access key"
  type        = string
  sensitive   = true
}

variable "target_secret_key" {
  description = "Target account secret key"
  type        = string
  sensitive   = true
}


