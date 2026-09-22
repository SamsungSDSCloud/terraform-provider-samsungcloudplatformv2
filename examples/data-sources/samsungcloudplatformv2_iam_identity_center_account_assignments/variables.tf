variable "instance_id" {
  type        = string
  description = "The ID of the Identity Center Instance"
  default     = "ENTER YOUR RESOURCE'S INSTANCE_ID"
}

variable "target_account_id" {
  type    = string
  default = "ENTER YOUR RESOURCE'S TARGET_ACCOUNT_ID"
}

variable "size" {
  type    = number
  default = 100
}

variable "page" {
  type    = number
  default = null
}

variable "sort" {
  type    = string
  default = null
}

variable "target_account_name" {
  type    = string
  default = null
}

variable "target_account_email" {
  type    = string
  default = null
}

variable "permission_set_id" {
  type    = string
  default = "ENTER YOUR RESOURCE'S PERMISSION_SET_ID"
}

variable "principal_name" {
  type    = string
  default = null
}

variable "role_srn" {
  type    = string
  default = null
}



