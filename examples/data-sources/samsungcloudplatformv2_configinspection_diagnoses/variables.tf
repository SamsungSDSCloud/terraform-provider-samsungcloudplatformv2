variable "with_count" {
  description = "With count"
  type        = string
  default     = null
}

variable "limit" {
  description = "Limit"
  type        = number
  default     = null
}

variable "marker" {
  description = "Marker"
  type        = string
  default     = null
}

variable "sort" {
  description = "Sort"
  type        = string
  default     = null
}

variable "account_id" {
  description = "Account Id"
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
  description = "Diagnosis state"
  type        = string
  default     = null
}

variable "start_date" {
  description = "Start date"
  type        = string
  default     = null
}

variable "end_date" {
  description = "End date"
  type        = string
  default     = null
}

variable "user_id" {
  description = "User id"
  type        = string
  default     = "ENTER YOUR RESOURCE'S USER_ID"
}



