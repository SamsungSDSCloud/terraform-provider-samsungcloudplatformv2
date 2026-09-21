variable "landing_zone_id" {
  description = "Landing Zone ID"
  type        = string
}

variable "name" {
  description = "Account Name"
  type        = string
}

variable "parent_unit_id" {
  description = "Parent Organization Unit ID"
  type        = string
}

variable "email" {
  description = "Account Email"
  type        = string
}

variable "sso_user_email" {
  description = "SSO User Email"
  type        = string
  default     = "test1010@samsung.com"
}

variable "sso_user_name" {
  description = "SSO User Name"
  type        = string
  default     = "test1010"
}

variable "sso_user_real_name" {
  description = "SSO User Real Name"
  type        = string
  default     = "test1010"
}


