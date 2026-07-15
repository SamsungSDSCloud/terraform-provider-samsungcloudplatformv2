variable "size" {
  type    = number
  default = 10
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
  description = "Filter by Internet Gateway ID"
  type        = string
  default     = "ENTER YOUR RESOURCE'S ID"
}

variable "name" {
  description = "Filter by Internet Gateway name"
  type        = string
  default     = null
}

variable "type" {
  description = "Filter by Internet Gateway type"
  type        = string
  default     = null
}

variable "state" {
  description = "Filter by Internet Gateway state"
  type        = string
  default     = null
}

variable "vpc_id" {
  description = "Filter by VPC ID"
  type        = string
  default     = "ENTER YOUR RESOURCE'S VPC_ID"
}

variable "vpc_name" {
  description = "Filter by VPC name"
  type        = string
  default     = null
}



