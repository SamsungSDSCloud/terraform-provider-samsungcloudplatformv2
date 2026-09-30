variable "workspace_name" {
  description = "Workspace name."
  type        = string
  default     = "tf-test-workspace"
}

variable "resource_id" {
  description = "Virtual server."
  type        = string
  default     = "ENTER YOUR RESOURCE'S RESOURCE_ID"
}

variable "resource_name" {
  description = "Resource name."
  type        = string
  default     = ""
}

variable "resource_type" {
  description = "Resource type: virtual_server."
  type        = string
  default     = "virtual_server"
}



