variable "size" {
  type    = number
  default = 0
}

variable "page" {
  type    = number
  default = 0
}

variable "sort" {
  type    = string
  default = ""
}

variable "id" {
  description = "Filter by Internet Gateway ID"
  type        = string
  default     = ""
}

variable "name" {
  description = "Filter by Internet Gateway name"
  type        = string
  default     = ""
}

variable "type" {
  description = "Filter by Internet Gateway type"
  type        = string
  default     = ""
}

variable "state" {
  description = "Filter by Internet Gateway state"
  type        = string
  default     = ""
}

variable "vpc_id" {
  description = "Filter by VPC ID"
  type        = string
  default     = ""
}

variable "vpc_name" {
  description = "Filter by VPC name"
  type        = string
  default     = ""
}


