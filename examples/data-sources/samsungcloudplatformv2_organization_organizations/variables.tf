variable "size" {
  description = "Number of results per page"
  type        = number
  default     = null
}

variable "page" {
  description = "Page number to retrieve (0-based)"
  type        = number
  default     = null
}

variable "sort" {
  description = "Sort criteria in the format 'field:direction' (e.g., 'created_at:desc', 'id:asc')"
  type        = string
  default     = null

  validation {
    condition     = var.sort == null || can(regex("^[a-z_]+:(asc|desc)$", var.sort))
    error_message = "sort must be in the format 'field:asc' or 'field:desc' (e.g., 'created_at:desc')."
  }
}

variable "name" {
  description = "Filter organizations by name"
  type        = string
  default     = null
}

variable "master_account_id" {
  description = "Filter organizations by master account ID"
  type        = string
  default     = "ENTER YOUR RESOURCE'S MASTER_ACCOUNT_ID"
}

variable "organization_id" {
  description = "Filter by specific organization ID"
  type        = string
  default     = "ENTER YOUR RESOURCE'S ORGANIZATION_ID"
}


