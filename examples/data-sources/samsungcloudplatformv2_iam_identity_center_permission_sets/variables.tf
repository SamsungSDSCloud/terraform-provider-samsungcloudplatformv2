variable "instance_id" {
  type        = string
  description = "The ID of the Identity Center Instance"
  default     = "ENTER YOUR RESOURCE'S INSTANCE_ID"
}

variable "name" {
  type    = string
  default = null
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



