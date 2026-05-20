variable "size" {
  type        = number
  description = "Number of items to return per page (minimum: 0)"
  default     = 0
}

variable "page" {
  type        = number
  description = "Page number (minimum: 0)"
  default     = 0
}

variable "sort" {
  type        = string
  description = "Sort order (e.g., created_at:desc)"
  default     = ""
}

variable "ip_address" {
  type        = string
  description = "Filter by IP address"
  default     = ""
}

variable "state" {
  type        = string
  description = "Filter by PublicIP state (RESERVED | ATTACHED | DELETED)"
  default     = ""
}

variable "attached_resource_type" {
  type        = string
  description = "Filter by attached resource type (VM | ALB | LB | BM | DB | NAT_GW | GPU_NODE | VPN | GPU_SERVER | EPAS | POSTGRESQL | MARIADB | SQLSERVER | CACHESTORE | SCALABLEDB | EVENTSTREAMS | SEARCHENGINE | VERTICA | SUBNET | MYSQL)"
  default     = ""
}

variable "attached_resource_id" {
  type        = string
  description = "Filter by attached resource ID"
  default     = ""
}

variable "attached_resource_name" {
  type        = string
  description = "Filter by attached resource name"
  default     = ""
}

variable "vpc_id" {
  type        = string
  description = "Filter by VPC ID"
  default     = ""
}

variable "type" {
  type        = string
  description = "Filter by PublicIP type (IGW | GGW | SIGW)"
  default     = ""
}


