variable "instance_id" {
  type        = string
  description = "IAM Identity Center Instance ID"
  default     = "ENTER YOUR RESOURCE'S INSTANCE_ID"
}

variable "name" {
  type        = string
  description = "Permission Set Name"
  default     = "Admin_Permission_Set1"
}

variable "description" {
  type        = string
  description = "Permission Set Description"
  default     = "Permission set for administrator"
}

variable "session_duration" {
  type        = number
  description = "Maximum Session Duration in seconds (900-43200)"
  default     = 3600
}



