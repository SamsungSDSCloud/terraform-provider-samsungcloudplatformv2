variable "notification_group_name" {
  description = "Notification group name."
  type        = string
  default     = "tf-test-notification-group"
}

variable "description" {
  description = "Notification group description."
  type        = string
  default     = "Terraform test notification group"
}

variable "recipient_user_ids" {
  description = "List of recipient user IDs."
  type        = list(string)
  default     = []
}



