variable "budget_id" {
  type        = string
  default     = "ENTER YOUR RESOURCE'S BUDGET_ID"
  description = "Budget ID (optional)"
}

variable "budget_name" {
  type        = string
  default     = null
  description = "Budget name (optional)"
}
variable "search_name" {
  type        = string
  default     = null
  description = "Seach name (optional)"
}

variable "sort" {
  type        = string
  default     = null
  description = "sort (optional)"
}


variable "page" {
  type        = number
  default     = 0
  description = "page (optional)"
}

variable "size" {
  type        = number
  default     = 20
  description = "count (optional)"
}


