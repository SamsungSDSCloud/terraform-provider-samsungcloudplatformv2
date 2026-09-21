variable "instance_id" {
  type        = string
  description = "The ID of the Identity Center Instance"
  default     = "ENTER YOUR RESOURCE'S INSTANCE_ID"
}

variable "name" {
  type        = string
  description = "The name of the IAM Identity Center Group"
  default     = "Admin_Groupss"
}

variable "description" {
  type        = string
  description = "The description of the IAM Identity Center Group"
  default     = "Group for administrators"
}



