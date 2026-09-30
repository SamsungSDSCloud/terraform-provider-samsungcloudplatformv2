variable "config_data" {
  description = "Base64 encoded YAML alert manager configuration file."
  type        = string
  default     = "Z2xvYmFsOgogIHJlc29sdmVfdGltZW91dDogNW0Kcm91dGU6CiAgcmVjZWl2ZXI6IGRlZmF1bHQK"
}

variable "notification_group_id" {
  description = "Notification group ID."
  type        = string
  default     = "ENTER YOUR RESOURCE'S NOTIFICATION_GROUP_ID"
}

variable "workspace_id" {
  description = "Workspace ID."
  type        = string
  default     = "ENTER YOUR RESOURCE'S WORKSPACE_ID"
}



