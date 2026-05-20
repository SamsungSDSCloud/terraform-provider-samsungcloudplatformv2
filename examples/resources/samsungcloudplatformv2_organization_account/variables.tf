variable "login_id" {
  description = "Login ID"
  type        = string
  default     = ""
}

variable "name" {
  description = "Account Name"
  type        = string
  default     = ""
}

variable "organization_id" {
  description = "Organization ID"
  type        = string
  default     = ""
}

variable "role_name" {
  description = "Role Name"
  type        = string
  default     = ""
}

variable "lazy_policy" {
  description = "Linked Policy Query YN"
  type        = bool
  default     = false
}

variable "parent_unit_id" {
  description = "Parent Organization Unit ID (for update/move)"
  type        = string
  default     = ""
}

variable "account_id" {
  description = "Account ID (for delete)"
  type        = string
  default     = ""
}

