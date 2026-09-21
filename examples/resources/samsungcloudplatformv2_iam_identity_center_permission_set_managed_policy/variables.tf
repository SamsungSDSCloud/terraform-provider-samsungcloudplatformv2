variable "instance_id" {
  type        = string
  description = "The ID of the Identity Center Instance"
  default     = "ENTER YOUR RESOURCE'S INSTANCE_ID"
}

variable "permission_set_id" {
  type        = string
  description = "The ID of the IAM Identity Center Permission Set"
  default     = "ENTER YOUR RESOURCE'S PERMISSION_SET_ID"
}

variable "managed_policy_id" {
  type        = string
  description = "The ID of the IAM Managed Policy"
  default     = "ENTER YOUR RESOURCE'S MANAGED_POLICY_ID"
}

variable "managed_policy_name" {
  type        = string
  description = "IAM Managed Policy Display Name"
  default     = "AdministratorAccess"
}



