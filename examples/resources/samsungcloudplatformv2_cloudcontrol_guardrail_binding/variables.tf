variable "guardrail_ids" {
  description = "List of guardrail IDs to enable"
  type        = list(string)
  default     = ["ENTER YOUR RESOURCE'S GUARDRAIL_IDS"]
}

variable "unit_ids" {
  description = "List of organization unit IDs to apply guardrails to"
  type        = list(string)
  default     = ["ENTER YOUR RESOURCE'S UNIT_IDS"]
}

variable "landing_zone_id" {
  description = "Landing Zone ID that contains the organization units"
  type        = string
  default     = "ENTER YOUR RESOURCE'S LANDING_ZONE_ID"
}


