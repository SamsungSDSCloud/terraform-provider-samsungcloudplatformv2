variable "name" {
  type        = string
  default     = "testregistry"
  description = "The name of the container registry."
}

variable "public_visible_enabled" {
  type        = bool
  default     = false
  description = "Whether the registry is publicly visible."
}

variable "public_endpoint_enabled" {
  type        = bool
  default     = false
  description = "Whether the public endpoint is enabled."
}

variable "private_acl_enabled" {
  type        = bool
  default     = false
  description = "Whether private ACL is enabled."
}

variable "public_acl_enabled" {
  type        = bool
  default     = false
  description = "Whether public ACL is enabled."
}

variable "private_acl_resources" {
  type = list(object({
    resource_id   = string
    resource_name = string
    resource_type = string
    resource_ips  = list(string)
  }))
  default     = []
  description = "List of resources allowed private access to the registry."
}

variable "public_acl_resources" {
  type = list(object({
    resource_id   = string
    resource_name = string
    resource_type = string
    resource_ips  = list(string)
  }))
  default     = []
  description = "List of resources allowed public access to the registry."
}



