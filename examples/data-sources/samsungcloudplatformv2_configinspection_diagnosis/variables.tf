variable "diagnosis_id" {
  description = "Id of diagnosis"
  type        = string
  default     = "ENTER YOUR RESOURCE'S DIAGNOSIS_ID"
}

variable "diagnosis_request_sequence" {
  description = "Sequence of diagnosis request"
  type        = string
  default     = "SCPCIS-EC862048D7744453Axxxxxxxxxxxxxxx"
}

variable "with_count" {
  description = "Whether to include the total item count in the response"
  type        = string
  default     = true
}

variable "limit" {
  description = "Maximum number of items to return per page"
  type        = number
  default     = 5
}

variable "marker" {
  description = "Pagination token from a previous response to fetch the next page"
  type        = string
  default     = null
}

variable "sort" {
  description = "The sorting criteria in the format 'field_name:asc' for ascending or 'field_name:desc' for descending order"
  type        = string
  default     = null
}



