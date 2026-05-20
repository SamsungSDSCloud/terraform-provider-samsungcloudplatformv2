variable "name" {
  type        = string
  description = "Organization Unit name"
  default     = ""
}

variable "description" {
  type        = string
  description = "Organization Unit description"
  default     = ""
}

variable "organization_id" {
  type        = string
  description = "Organization ID"
  default     = ""
}

variable "parent_unit_id" {
  type        = string
  description = "Parent Organization Unit ID (required)"
  default     = ""
}

variable "policy_ids" {
  type        = list(string)
  description = "Policy IDs"
  default     = [""]
}

