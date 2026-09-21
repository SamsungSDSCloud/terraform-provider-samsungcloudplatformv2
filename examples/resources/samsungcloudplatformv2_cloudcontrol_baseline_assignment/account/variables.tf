variable "assignment_id" {
  description = "Root/Organization Unit/Account ID to assign baseline to"
  type        = string
}

variable "landing_zone_id" {
  description = "Landing Zone ID that contains this baseline assignment"
  type        = string
}

variable "resource_type" {
  description = "Type of resource to assign baseline to (ACCOUNT or OU)"
  type        = string
}

variable "parent_unit_id" {
  description = "Parent Organization Unit ID"
  type        = string
  default     = "ENTER YOUR RESOURCE'S PARENT_UNIT_ID"
}

variable "sso_user_name" {
  description = "SSO User Name for the assignment"
  type        = string
  default     = ""
}

variable "sso_user_real_name" {
  description = "Real name of the SSO User"
  type        = string
  default     = ""
}



