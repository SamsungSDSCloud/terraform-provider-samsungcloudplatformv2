variable "diagnosis_id" {
  description = "Id of diagnosis"
  type        = string
  default     = "ENTER YOUR RESOURCE'S DIAGNOSIS_ID"
}

variable "diagnosis_request_sequence" {
  description = "Sequence of diagnosis request"
  type        = string
  default     = "SCPCIS-EC862048D7744453A759493BCC165BAC"
}

variable "with_count" {
  description = "With count"
  type        = string
  default     = true
}

variable "limit" {
  description = "Limit"
  type        = number
  default     = 5
}

variable "marker" {
  description = "Marker"
  type        = string
  default     = null
}

variable "sort" {
  description = "Sort"
  type        = string
  default     = null
}



