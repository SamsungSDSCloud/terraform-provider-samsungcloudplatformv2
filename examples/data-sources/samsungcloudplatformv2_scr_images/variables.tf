variable "repository_id" {
  type        = string
  default     = "ENTER YOUR RESOURCE'S REPOSITORY_ID"
  description = "The ID of the repository to list images from."
}

variable "name" {
  type        = string
  default     = ""
  description = "Filter images by name."
}

variable "sort" {
  type        = string
  default     = ""
  description = "Sort order, e.g. name:asc."
}

variable "page" {
  type        = number
  default     = 0
  description = "Page number (0-based)."
}

variable "size" {
  type        = number
  default     = 20
  description = "Number of items per page."
}



