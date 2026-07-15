variable "size" {
  type        = number
  description = "Size"
  default     = 20
}

variable "page" {
  type        = number
  description = "Page"
  default     = 0
}

variable "sort" {
  type        = string
  description = "Sort"
  default     = "created_at:desc"
}

variable "name" {
  type        = string
  description = "NAT Gateway Name"
  default     = null
}

variable "vpc_id" {
  type        = string
  description = "VPC ID"
  default     = "ENTER YOUR RESOURCE'S VPC_ID"
}

variable "vpc_name" {
  type        = string
  description = "VPC Name"
  default     = null
}

variable "subnet_id" {
  type        = string
  description = "Subnet ID"
  default     = "ENTER YOUR RESOURCE'S SUBNET_ID"
}

variable "subnet_name" {
  type        = string
  description = "Subnet Name"
  default     = null
}

variable "nat_gateway_ip_address" {
  type        = string
  description = "NAT Gateway IP Address"
  default     = null
}

variable "state" {
  type        = string
  description = "NAT Gateway State"
  default     = null
}



