variable "security_group_id" {
  type    = string
  default = "ENTER YOUR RESOURCE'S SECURITY_GROUP_ID"
}

variable "size" {
  type    = number
  default = null
}

variable "page" {
  type    = number
  default = 0
}

variable "id" {
  type    = string
  default = "ENTER YOUR RESOURCE'S ID"
}

variable "remote_ip_prefix" {
  type    = string
  default = null
}

variable "remote_group_id" {
  type    = string
  default = "ENTER YOUR RESOURCE'S REMOTE_GROUP_ID"
}

variable "description" {
  type    = string
  default = null
}

variable "direction" {
  type    = string
  default = null
}

variable "service" {
  type    = string
  default = null
}


