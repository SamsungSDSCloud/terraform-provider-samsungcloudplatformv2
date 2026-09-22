variable "instance_id" {
  type        = string
  description = "IAM Identity Center Instance ID"
  default     = "ENTER YOUR RESOURCE'S INSTANCE_ID"
}

variable "target_account_id" {
  type        = string
  description = "Target account ID"
  default     = "ENTER YOUR RESOURCE'S TARGET_ACCOUNT_ID"
}

variable "principal_id" {
  type        = string
  description = "Principal ID (user or group)"
  default     = "ENTER YOUR RESOURCE'S PRINCIPAL_ID"
}

variable "principal_type" {
  type        = string
  description = "Principal type: USER or GROUP"
  default     = "USER"
}

variable "permission_set_id" {
  type        = string
  description = "Permission set ID"
  default     = "ENTER YOUR RESOURCE'S PERMISSION_SET_ID"
}



