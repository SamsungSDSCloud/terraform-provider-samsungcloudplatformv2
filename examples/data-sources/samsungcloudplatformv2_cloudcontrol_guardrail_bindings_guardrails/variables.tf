variable "target_id" {
  description = "Target ID to query guardrails for"
  type        = string
}
variable "landing_zone_id" {
  description = "Landing Zone ID that contains the organization unit"
  type        = string
  default     = "ENTER YOUR RESOURCE'S LANDING_ZONE_ID"
}
variable "name" {
  description = "Filter by guardrail name"
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


