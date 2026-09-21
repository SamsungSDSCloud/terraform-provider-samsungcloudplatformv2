variable "instance_id" {
  type        = string
  description = "The ID of the Identity Center Instance"
  default     = "ENTER YOUR RESOURCE'S INSTANCE_ID"
}

variable "user_id" {
  type    = string
  default = "ENTER YOUR RESOURCE'S USER_ID"
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

variable "excluded_group_id" {
  type    = string
  default = "ENTER YOUR RESOURCE'S EXCLUDED_GROUP_ID"
}

variable "excluded_account_id" {
  type    = string
  default = "ENTER YOUR RESOURCE'S EXCLUDED_ACCOUNT_ID"
}



