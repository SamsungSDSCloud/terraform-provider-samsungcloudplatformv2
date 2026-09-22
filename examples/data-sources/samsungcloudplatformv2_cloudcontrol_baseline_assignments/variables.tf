variable "assignment_id" {
  description = "Root/Organization Unit/Account ID to filter by"
  type        = string
  default     = "ENTER YOUR RESOURCE'S ASSIGNMENT_ID"
}

variable "landing_zone_id" {
  description = "Landing Zone ID to filter by"
  type        = string
  default     = "ENTER YOUR RESOURCE'S LANDING_ZONE_ID"
}

variable "resource_type" {
  description = "Resource type to filter by (OU or ACCOUNT)"
  type        = string
  default     = ""
}

variable "status" {
  description = "Baseline assignment status to filter by"
  type        = string
  default     = ""
}


