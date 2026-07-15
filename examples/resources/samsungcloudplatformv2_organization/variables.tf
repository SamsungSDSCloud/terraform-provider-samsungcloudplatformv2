variable "name" {
  type        = string
  description = "Organization name"
  default     = "terraform-org"
}

variable "delegation_account_id" {
  type        = string
  description = "Delegation Account ID"
  default     = "ENTER YOUR RESOURCE'S DELEGATION_ACCOUNT_ID"
}

variable "use_scp_yn" {
  type        = bool
  description = "Use SCP"
  default     = true
}


