variable "size" {
  type        = number
  description = "size"
  default     = 0
}

variable "page" {
  type        = number
  description = "page"
  default     = 0
}

variable "sort" {
  type        = string
  description = "sort"
  default     = ""
}

variable "name" {
  type        = string
  description = "Private NAT Name"
  default     = ""
}

variable "cidr" {
  type        = string
  description = "Private NAT IP range"
  default     = ""
}

variable "vpc_id" {
  type        = string
  description = "VPC Id"
  default     = ""
}

variable "service_resource_id" {
  type        = string
  description = "Private NAT connected Service Resource ID"
  default     = ""
}

variable "service_type" {
  type        = string
  description = "Private NAT connected Service Type"
  default     = ""
}

variable "service_resource_name" {
  type        = string
  description = "Private NAT connected Service Resource Name"
  default     = ""
}

variable "state" {
  type        = string
  description = "Private NAT State"
  default     = ""
}


