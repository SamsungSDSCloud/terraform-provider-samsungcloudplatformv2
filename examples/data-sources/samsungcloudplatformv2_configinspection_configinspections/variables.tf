variable "with_count" {
  description = "Whether to include the total item count in the response"
  type        = string
  default     = null // true, false
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
  description = "Sort results as 'field:asc' or 'field:desc'"
  type        = string
  default     = null
}

variable "is_mine" {
  description = "My Config Inspection"
  type        = bool
  default     = null
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

variable "diagnosis_account_id" {
  description = "Account Id of diagnosis"
  type        = string
  default     = "ENTER YOUR RESOURCE'S DIAGNOSIS_ACCOUNT_ID"
}

variable "recent_diagnosis_state" {
  description = "Filter by the latest diagnosis status"
  type        = list(string)
  default     = null
}

variable "start_date" {
  description = "Include only inspections created on or after this date"
  type        = string
  default     = null
}

variable "end_date" {
  description = "Include only inspections created on or before this date"
  type        = string
  default     = null
}



