variable "guardrail_id" {
  description = "Guardrail ID to query targets for"
  type        = string
}

variable "target_type" {
  description = "Target type to filter by (ACCOUNT or OU)"
  type        = string
}

variable "landing_zone_id" {
  description = "Landing Zone ID"
  type        = string
  default     = "ENTER YOUR RESOURCE'S LANDING_ZONE_ID"
}

variable "name" {
  description = "Filter by target name"
  type        = string
  default     = ""
}

variable "size" {
  description = "Number of items to return"
  type        = number
  default     = 0
}

variable "page" {
  description = "Page number"
  type        = number
  default     = 0
}

variable "sort" {
  description = "Sort criteria"
  type        = string
  default     = ""
}


