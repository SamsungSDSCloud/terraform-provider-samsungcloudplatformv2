variable "name" {
  type        = string
  description = "Organization name"
  default     = ""
}

variable "delegation_account_id" {
  type        = string
  description = "Delegation Account ID"
  default     = ""
}

variable "use_scp_yn" {
  type        = bool
  description = "Use SCP"
  default     = false
}

