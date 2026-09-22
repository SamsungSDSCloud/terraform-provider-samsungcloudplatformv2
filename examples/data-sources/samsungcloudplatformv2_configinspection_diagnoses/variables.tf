variable "with_count" {
  description = "Whether to include the total item count in the response"
  type        = string
  default     = null
}

variable "limit" {
  description = "Maximum number of items to return per page"
  type        = number
  default     = null
}

variable "marker" {
  description = "Pagination token from a previous response to fetch the next page"
  type        = string
  default     = null
}

variable "sort" {
  description = "The sorting criteria in the format 'field_name:asc' for ascending or 'field_name:desc' for descending order"
  type        = string
  default     = null
}

variable "account_id" {
  description = "Account Identifier"
  type        = string
  default     = "ENTER YOUR RESOURCE'S ACCOUNT_ID"
}

variable "diagnosis_id" {
  description = "Id of diagnosis"
  type        = string
  default     = "ENTER YOUR RESOURCE'S DIAGNOSIS_ID"
}

variable "diagnosis_name" {
  description = "Name of diagnosis"
  type        = string
  default     = null
}

variable "csp_type" {
  description = "Type of cloud service provider"
  type        = string
  default     = null
}

variable "diagnosis_state" {
  description = "Status of diagnosis"
  type        = string
  default     = null
}

variable "start_date" {
  description = "Include only items created on or after this date"
  type        = string
  default     = null
}

variable "end_date" {
  description = "Include only items created on or before this date"
  type        = string
  default     = null
}

variable "user_id" {
  description = "Account owner of this diagnosis"
  type        = string
  default     = "ENTER YOUR RESOURCE'S USER_ID"
}



