variable "size" {
  type        = number
  description = "size"
  default     = 5
}

variable "page" {
  type        = number
  description = "page"
  default     = 0
}

variable "sort" {
  type        = string
  description = "sort"
  default     = "created_at:desc"
}

variable "name" {
  type        = string
  description = "Private NAT Name"
  default     = null
}

variable "cidr" {
  type        = string
  description = "Private NAT IP range"
  default     = null
}

variable "vpc_id" {
  type        = string
  description = "VPC Id"
  default     = "ENTER YOUR RESOURCE'S VPC_ID"
}

variable "service_resource_id" {
  type        = string
  description = "Private NAT connected Service Resource ID"
  default     = "ENTER YOUR RESOURCE'S SERVICE_RESOURCE_ID"
}

variable "service_type" {
  type        = string
  description = "Private NAT connected Service Type"
  default     = null
}

variable "service_resource_name" {
  type        = string
  description = "Private NAT connected Service Resource Name"
  default     = null
}

variable "state" {
  type        = string
  description = "Private NAT State"
  default     = null
}



