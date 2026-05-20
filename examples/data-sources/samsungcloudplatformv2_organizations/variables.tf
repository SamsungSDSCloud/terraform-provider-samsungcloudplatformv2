variable "size" {
  description = "Number of results per page"
  type        = number
  default     = 0
}

variable "page" {
  description = "Page number to retrieve (0-based)"
  type        = number
  default     = 0
}

variable "sort" {
  description = "Sort criteria in the format 'field:direction' (e.g., 'created_at:desc', 'id:asc')"
  type        = string

  validation {
    condition     = var.sort == null || can(regex("^[a-z_]+:(asc|desc)$", var.sort))
    error_message = "sort must be in the format 'field:asc' or 'field:desc' (e.g., 'created_at:desc')."
  }
  default = ""
}

variable "name" {
  description = "Filter organizations by name"
  type        = string
  default     = ""
}

variable "master_account_id" {
  description = "Filter organizations by master account ID"
  type        = string
  default     = ""
}

variable "organization_id" {
  description = "Filter by specific organization ID"
  type        = string
  default     = ""
}

