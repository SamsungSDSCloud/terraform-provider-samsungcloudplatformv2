variable "cluster_id" {
  type        = string
  default     = "ENTER YOUR RESOURCE'S CLUSTER_ID"
  description = "ID of the SKE cluster to manage linked resources for"
}

variable "linked_resource_id" {
  type        = string
  default     = "ENTER YOUR RESOURCE'S LINKED_RESOURCE_ID"
  description = "ID of the linked resource"
}

variable "linked_resource_name" {
  type        = string
  default     = "vulinh-object-storage"
  description = "Name of the linked resource"
}

variable "linked_resource_type" {
  type        = string
  default     = "obs"
  description = "Type of the linked resource (fs or obs)"
}



