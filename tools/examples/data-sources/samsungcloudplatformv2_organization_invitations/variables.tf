variable "organization_id" {
  type        = string
  description = "Filter by Organization ID"
  default     = "ENTER YOUR RESOURCE'S ORGANIZATION_ID"
}

variable "size" {
  type        = number
  description = "Number of results per page"
  default     = null
}

variable "page" {
  type        = number
  description = "Page number"
  default     = null
}

variable "sort" {
  type        = string
  description = "Sort criteria (e.g., 'created_at:desc')"
  default     = null
}

variable "account_id" {
  type        = string
  description = "Filter by Account ID"
  default     = "ENTER YOUR RESOURCE'S ACCOUNT_ID"
}

variable "account_name" {
  type        = string
  description = "Filter by Account Name"
  default     = null
}

variable "account_email" {
  type        = string
  description = "Filter by Account Email"
  default     = null
}

variable "state" {
  type        = string
  description = "Filter by Invitation State"
  default     = null
}

variable "login_id" {
  type        = string
  description = "Filter by Login ID"
  default     = "ENTER YOUR RESOURCE'S LOGIN_ID"
}


