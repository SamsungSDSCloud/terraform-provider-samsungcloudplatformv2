variable "size" {
  type        = number
  description = "Size"
  default     = 0
}

variable "page" {
  type        = number
  description = "Page"
  default     = 0
}

variable "sort" {
  type        = string
  description = "Sort"
  default     = ""
}

variable "name" {
  type        = string
  description = "NAT Gateway Name"
  default     = ""
}

variable "vpc_id" {
  type        = string
  description = "VPC ID"
  default     = ""
}

variable "vpc_name" {
  type        = string
  description = "VPC Name"
  default     = ""
}

variable "subnet_id" {
  type        = string
  description = "Subnet ID"
  default     = ""
}

variable "subnet_name" {
  type        = string
  description = "Subnet Name"
  default     = ""
}

variable "nat_gateway_ip_address" {
  type        = string
  description = "NAT Gateway IP Address"
  default     = ""
}

variable "state" {
  type        = string
  description = "NAT Gateway State"
  default     = ""
}


