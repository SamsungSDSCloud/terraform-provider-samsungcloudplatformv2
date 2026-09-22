variable "landing_zone_id" {
  description = "Landing Zone ID to filter by"
  type        = string
  default     = "ENTER YOUR RESOURCE'S LANDING_ZONE_ID"
}

variable "name" {
  description = "Guardrail name to filter by"
  type        = string
  default     = ""
}

variable "guidance" {
  description = "Policy guidance to filter by (REQUIRED, ALARM, DISABLED)"
  type        = string
  default     = ""
}

variable "service_name" {
  description = "Service name to filter by"
  type        = string
  default     = ""
}

variable "status" {
  description = "Guardrail status to filter by (ENABLED, DISABLED)"
  type        = string
  default     = ""
}

variable "size" {
  description = "Page size to request"
  type        = number
  default     = 0
}

variable "page" {
  description = "Page number to request"
  type        = number
  default     = 0
}

variable "exclude_unit_id" {
  description = "Linked(Exclusion) Unit ID to filter by"
  type        = string
  default     = "ENTER YOUR RESOURCE'S EXCLUDE_UNIT_ID"
}
variable "sort" {
  description = "Sort criteria for the request, formatted as 'field:direction'"
  type        = string
  default     = ""
}


