variable "organization_id" {
  type        = string
  description = "Organization ID"
  default     = "ENTER YOUR RESOURCE'S ORGANIZATION_ID"
}

variable "service_type" {
  type        = string
  default     = ""
  description = "Service Type (optional)"
}

variable "account_id" {
  type        = string
  default     = "ENTER YOUR RESOURCE'S ACCOUNT_ID"
  description = "Delegation Account ID (optional)"
}

variable "size" {
  type        = number
  default     = 20
  description = "Page size (optional)"
}

variable "page" {
  type        = number
  default     = 0
  description = "Page number (optional)"
}

variable "sort" {
  type        = string
  default     = null
  description = "Sort criteria (optional)"
}


