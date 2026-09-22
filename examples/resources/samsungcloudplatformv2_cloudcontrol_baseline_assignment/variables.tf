variable "assignment_id" {
  description = "Root/Organization Unit/Account ID to assign baseline to"
  type        = string
}

variable "landing_zone_id" {
  description = "Landing Zone ID that contains this baseline assignment"
  type        = string
}

variable "resource_type" {
  description = "Type of resource to assign baseline to (ACCOUNT or OU)"
  type        = string
}

variable "agree_yn" {
  description = "Agreement to baseline assignment (Y or N)"
  type        = string
  default     = "Y"
}

variable "reregister_trigger" {
  description = "Arbitrary value used to force a re-registration (Update) of this baseline assignment. Only meaningful for resource_type = 'OU'."
  type        = string
  default     = "1"
}


