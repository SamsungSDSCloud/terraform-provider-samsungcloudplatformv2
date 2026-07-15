variable "login_id" {
  description = "Login ID"
  type        = string
  default     = "ENTER YOUR RESOURCE'S LOGIN_ID"
}

variable "name" {
  description = "Account Name"
  type        = string
  default     = "score"
}

variable "organization_id" {
  description = "Organization ID"
  type        = string
  default     = "ENTER YOUR RESOURCE'S ORGANIZATION_ID"
}

variable "role_name" {
  description = "Role Name"
  type        = string
  default     = "OrganizationAccountAccessRole"
}

variable "lazy_policy" {
  description = "Linked Policy Query YN"
  type        = bool
  default     = false
}

variable "parent_unit_id" {
  description = "Parent Organization Unit ID (for update/move)"
  type        = string
  default     = "ENTER YOUR RESOURCE'S PARENT_UNIT_ID"
}

variable "account_id" {
  description = "Account ID (for delete)"
  type        = string
  default     = "ENTER YOUR RESOURCE'S ACCOUNT_ID"
}


