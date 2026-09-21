variable "name" {
  type        = string
  description = "Organization name"
  default     = "terraform-org"
}

variable "use_scp_yn" {
  type        = bool
  description = "Use SCP"
  default     = true
}


