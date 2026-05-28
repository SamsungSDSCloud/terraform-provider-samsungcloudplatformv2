variable "size" {
  type        = number
  description = "Number of items to return per page (minimum: 0)"
  default     = 10
}

variable "page" {
  type        = number
  description = "Page number (minimum: 0)"
  default     = 0
}

variable "sort" {
  type        = string
  description = "Sort order (e.g., created_at:desc)"
  default     = null
}

variable "ip_address" {
  type        = string
  description = "Filter by IP address"
  default     = null
}

variable "state" {
  type        = string
  description = "Filter by PublicIP state (RESERVED | ATTACHED | DELETED)"
  default     = null
}

variable "attached_resource_type" {
  type        = string
  description = "Filter by attached resource type (VM | ALB | LB | BM | DB | NAT_GW | GPU_NODE | VPN | GPU_SERVER | EPAS | POSTGRESQL | MARIADB | SQLSERVER | CACHESTORE | SCALABLEDB | EVENTSTREAMS | SEARCHENGINE | VERTICA | SUBNET | MYSQL)"
  default     = null
}

variable "attached_resource_id" {
  type        = string
  description = "Filter by attached resource ID"
  default     = "ENTER YOUR RESOURCE'S ATTACHED_RESOURCE_ID"
}

variable "attached_resource_name" {
  type        = string
  description = "Filter by attached resource name"
  default     = null
}

variable "vpc_id" {
  type        = string
  description = "Filter by VPC ID"
  default     = "ENTER YOUR RESOURCE'S VPC_ID"
}

variable "type" {
  type        = string
  description = "Filter by PublicIP type (IGW | GGW | SIGW)"
  default     = null
}



