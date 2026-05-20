variable "organization_id" {
  type        = string
  description = "Filter by Organization ID"
  default     = ""
}

variable "size" {
  type        = number
  description = "Number of results per page"
  default     = 0
}

variable "page" {
  type        = number
  description = "Page number"
  default     = 0
}

variable "sort" {
  type        = string
  description = "Sort criteria (e.g., 'created_at:desc')"
  default     = ""
}

variable "account_id" {
  type        = string
  description = "Filter by Account ID"
  default     = ""
}

variable "account_name" {
  type        = string
  description = "Filter by Account Name"
  default     = ""
}

variable "account_email" {
  type        = string
  description = "Filter by Account Email"
  default     = ""
}

variable "state" {
  type        = string
  description = "Filter by Invitation State"
  default     = ""
}

variable "login_id" {
  type        = string
  description = "Filter by Login ID"
  default     = ""
}

