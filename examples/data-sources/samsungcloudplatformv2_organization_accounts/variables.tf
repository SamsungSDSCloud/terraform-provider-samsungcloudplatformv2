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

variable "id" {
  type        = string
  description = "Filter by Account ID"
  default     = ""
}

variable "name" {
  type        = string
  description = "Filter by Account Name"
  default     = ""
}

variable "email" {
  type        = string
  description = "Filter by Account Email"
  default     = ""
}

variable "login_id" {
  type        = string
  description = "Filter by Login ID"
  default     = ""
}

variable "joined_start_date" {
  type        = string
  description = "Filter by joined start date (e.g., '2026-04-11T12:12:12.123Z')"
  default     = ""
}

variable "joined_end_date" {
  type        = string
  description = "Filter by joined end date (e.g., '2026-04-11T12:12:12.123Z')"
  default     = ""
}

variable "joined_method" {
  type        = string
  description = "Filter by joined method"
  default     = ""
}

variable "exclude_policy_id" {
  type        = string
  description = "Filter by exclude policy ID"
  default     = ""
}

variable "parent_unit_id" {
  type        = string
  description = "Filter by parent unit ID"
  default     = ""
}

variable "parent_unit_name" {
  type        = string
  description = "Filter by parent unit name"
  default     = ""
}

variable "type" {
  type        = string
  description = "Filter by account type"
  default     = ""
}

