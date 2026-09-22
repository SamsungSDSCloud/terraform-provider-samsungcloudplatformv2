variable "policy_id" {
  type        = string
  description = "Control policy ID whose assigned targets are listed"
}

variable "target_type" {
  type        = string
  description = "Target type (ROOT, ACCOUNT, OU)"
}

variable "organization_id" {
  type        = string
  description = "Organization ID"
  default     = "ENTER YOUR RESOURCE'S ORGANIZATION_ID"
}

variable "policy_category" {
  type        = string
  description = "Policy category filter (SCP, RCP, TAG, DGP)"
  default     = null
}

variable "name" {
  type        = string
  description = "Filter by target name"
  default     = null
}

variable "size" {
  type        = number
  description = "Number of results per page"
  default     = 20
}

variable "page" {
  type        = number
  description = "Page number"
  default     = null
}

variable "sort" {
  type        = string
  description = "Sort criteria (e.g., 'created_at:desc')"
  default     = null
}



