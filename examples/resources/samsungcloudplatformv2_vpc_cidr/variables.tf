variable "vpc_id" {
  description = "VPC ID to add CIDR"
  type        = string
  default     = ""
}

variable "cidr" {
  description = "CIDR block to add to the VPC (e.g., 192.168.0.0/24)"
  type        = string
  default     = ""
}


