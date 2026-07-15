variable "vpc_id" {
  description = "VPC ID to add CIDR"
  type        = string
  default     = "ENTER YOUR RESOURCE'S VPC_ID"
}

variable "cidr" {
  description = "CIDR block to add to the VPC (e.g., 192.168.0.0/24)"
  type        = string
  default     = "192.169.0.0/16"
}



