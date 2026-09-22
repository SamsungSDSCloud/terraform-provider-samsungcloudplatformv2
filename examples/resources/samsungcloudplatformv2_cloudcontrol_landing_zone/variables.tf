variable "additional_ou_name" {
  description = "Additional Organization Unit Name"
  type        = string
  default     = "SandboxOU"
}

variable "agree_yn" {
  description = "Terms Agreement YN"
  type        = string
  default     = "Y"
}

variable "audit_account_name" {
  description = "Audit Account Name"
  type        = string
  default     = "audit-account"
}

variable "audit_login_id" {
  description = "Audit Login ID"
  type        = string
  default     = "ENTER YOUR RESOURCE'S AUDIT_LOGIN_ID"
}

variable "basic_ou_name" {
  description = "Basic Organization Unit Name"
  type        = string
  default     = "SecurityOU"
}

variable "detective_guardrail_status" {
  description = "Detective Guardrail Status"
  type        = string
  default     = "ENABLED"
}

variable "log_archive_account_name" {
  description = "Log Archive Account Name"
  type        = string
  default     = "log-archive-account"
}

variable "log_archive_login_id" {
  description = "Log Archive Login ID"
  type        = string
  default     = "ENTER YOUR RESOURCE'S LOG_ARCHIVE_LOGIN_ID"
}

variable "sso_type" {
  description = "SSO Type"
  type        = string
  default     = "SELF"
}


