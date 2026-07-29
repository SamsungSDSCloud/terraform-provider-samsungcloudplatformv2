variable "organization_id" {
  type        = string
  default     = "ENTER YOUR RESOURCE'S ORGANIZATION_ID"
  description = "Organization ID"
}

variable "size" {
  type    = number
  default = 3
}

variable "page" {
  type    = number
  default = 0
}

variable "sort" {
  type    = string
  default = null
}

variable "id" {
  type    = string
  default = "ENTER YOUR RESOURCE'S ID"
}

variable "name" {
  type    = string
  default = null
}

variable "type" {
  type    = string
  default = null
}

variable "exclude_target_id" {
  type    = string
  default = "ENTER YOUR RESOURCE'S EXCLUDE_TARGET_ID"
}


