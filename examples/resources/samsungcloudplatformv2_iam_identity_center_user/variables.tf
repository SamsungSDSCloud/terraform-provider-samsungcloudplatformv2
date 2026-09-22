variable "instance_id" {
  type        = string
  description = "The ID of the Identity Center Instance"
  default     = "ENTER YOUR RESOURCE'S INSTANCE_ID"
}

variable "user_id" {
  type        = string
  description = "The user login ID for IAM Identity Center User"
  default     = "ENTER YOUR RESOURCE'S USER_ID"
}

variable "name" {
  type        = string
  description = "The real name of the IAM Identity Center User"
  default     = "John Doe4"
}

variable "email" {
  type        = string
  description = "The email of the IAM Identity Center User"
  default     = "john.doe2@example.com"
}

variable "password" {
  type        = string
  description = "The password for the IAM Identity Center User"
  default     = "ENTER YOUR RESOURCE'S PASSWORD"
  sensitive   = true
}



