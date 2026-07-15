variable "organization_id" {
  type        = string
  description = "Filter by Organization ID"
  default     = "ENTER YOUR RESOURCE'S ORGANIZATION_ID"
}

variable "size" {
  type        = number
  description = "Number of results per page"
  default     = 20
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

variable "id" {
  type        = string
  description = "Filter by Account ID"
  default     = "ENTER YOUR RESOURCE'S ID"
}

variable "name" {
  type        = string
  description = "Filter by Account Name"
  default     = null
}

variable "email" {
  type        = string
  description = "Filter by Account Email"
  default     = null
}

variable "login_id" {
  type        = string
  description = "Filter by Login ID"
  default     = "ENTER YOUR RESOURCE'S LOGIN_ID"
}

variable "joined_start_date" {
  type        = string
  description = "Filter by joined start date (e.g., '2026-04-11T12:12:12.123Z')"
  default     = null
}

variable "joined_end_date" {
  type        = string
  description = "Filter by joined end date (e.g., '2026-04-11T12:12:12.123Z')"
  default     = null
}

variable "joined_method" {
  type        = string
  description = "Filter by joined method"
  default     = null
}

variable "exclude_policy_id" {
  type        = string
  description = "Filter by exclude policy ID"
  default     = "ENTER YOUR RESOURCE'S EXCLUDE_POLICY_ID"
}

variable "parent_unit_id" {
  type        = string
  description = "Filter by parent unit ID"
  default     = "ENTER YOUR RESOURCE'S PARENT_UNIT_ID"
}

variable "parent_unit_name" {
  type        = string
  description = "Filter by parent unit name"
  default     = null
}

variable "type" {
  type        = string
  description = "Filter by account type"
  default     = null
}


